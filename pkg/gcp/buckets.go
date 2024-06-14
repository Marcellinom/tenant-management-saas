package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"google.golang.org/api/iterator"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BucketConfig struct {
	bucket, prefix string
}

// BucketBackend Hanya sebagai contoh jika mau implementasi interface custom backend,
//
//	terraform sudah support backend gcs dari sono nya
func BucketBackend(bucket, prefix string) *BucketConfig {
	if prefix == "" {
		log.Panic("sebaiknya prefix jangan kosong jendral")
	}
	return &BucketConfig{bucket: bucket, prefix: prefix}
}

// Init butuh tenant id dalam konteks
func (b *BucketConfig) Init(ct context.Context) error {
	tf, ok := ct.Value("terraform").(*tfexec.Terraform)
	if !ok {
		return fmt.Errorf("executable terraform tidak disediakan")
	}

	client, err := storage.NewClient(ct)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v\n", err)
	}
	defer client.Close()

	bucketName := b.bucket
	prefix := fmt.Sprintf("%s/%s", b.prefix, ct.Value("tenant_id").(string))
	localDir := tf.WorkingDir()

	// Get the bucket handle
	bucket := client.Bucket(bucketName)

	// Create a query with the specified prefix
	query := &storage.Query{Prefix: prefix}

	// Initialize an iterator to list objects
	it := bucket.Objects(ct, query)
	// Iterate through the objects
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to list objects: %v\n", err)
		}

		// Download the object
		if err = downloadObject(ct, client, bucketName, attrs.Name, localDir); err != nil {
			return fmt.Errorf("Failed to download object %s: %v\n", attrs.Name, err)
		}
	}
	err = tf.Init(ct)
	if err != nil {
		return err
	}
	return nil
}

func (b *BucketConfig) Apply(ctx context.Context) error {
	tf, ok := ctx.Value("terraform").(*tfexec.Terraform)
	if !ok {
		return fmt.Errorf("executable terraform tidak disediakan")
	}

	var tenant_env []tfexec.ApplyOption
	tenant_env, _ = ctx.Value("tenant_env").([]tfexec.ApplyOption)

	err := tf.Apply(ctx, tenant_env...)
	if err != nil {
		return err
	}

	// kirim state ke remote
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create upload client: %v\n", err)
	}
	bucketName := b.bucket
	remotePrefix := fmt.Sprintf("%s/%s", b.prefix, ctx.Value("tenant_id").(string))
	localDir := tf.WorkingDir()

	defer client.Close()
	return filepath.WalkDir(localDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(d.Name()) == ".tfstate" {
			relPath, err := filepath.Rel(localDir, path)
			if err != nil {
				return err
			}
			// Create the remote object name
			remotePath := filepath.Join(remotePrefix, relPath)
			remotePath = strings.ReplaceAll(remotePath, "\\", "/") // Ensure it's in Unix format
			return uploadFile(ctx, client, bucketName, path, remotePath)
		}
		return nil
	})
}

func uploadFile(ctx context.Context, client *storage.Client, bucketName, localPath, remotePath string) error {
	// Open the local file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file %s: %v", localPath, err)
	}
	defer file.Close()

	// Get a handle to the object
	wc := client.Bucket(bucketName).Object(remotePath).NewWriter(ctx)
	defer wc.Close()

	// Copy the file contents to the object
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("failed to copy file to bucket %s, object %s: %v", bucketName, remotePath, err)
	}

	fmt.Printf("Uploaded %s to gs://%s/%s\n", localPath, bucketName, remotePath)
	return nil
}

func downloadObject(ctx context.Context, client *storage.Client, bucketName, objectName, localDir string) error {
	// Get the object handle
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("new reader: %v", err)
	}
	defer rc.Close()

	// Create the local file
	localFilePath := localDir
	if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}

	potongan_prefix := strings.Split(objectName, "/")
	state_file := potongan_prefix[len(potongan_prefix)-1]

	// kalo ngebaca state remote dengan cara dipindahin ke local
	// di local namanya harus terraform.tfstate
	// gatau kenapa jangan tanya gw
	// mending pake built in terraform backend,
	// ini cuma contoh kalo bisa pake custom backend
	index := strings.Index(state_file, ".tfstate")
	if index != -1 {
		state_file = "terraform" + state_file[index:]
	}

	localFile, err := os.Create(filepath.Join(localDir, state_file))
	if err != nil {
		return fmt.Errorf("create file: %v", err)
	}
	defer localFile.Close()

	// Copy the object data to the local file
	if _, err := io.Copy(localFile, rc); err != nil {
		return fmt.Errorf("copy: %v", err)
	}

	fmt.Println(fmt.Printf("Downloaded %s to %s\n", objectName, localFilePath))
	return nil
}
