package gcp

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
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

func Bucket(bucket, prefix string) *BucketConfig {
	if prefix == "" {
		log.Panic("sebaiknya prefix jangan kosong jendral")
	}
	return &BucketConfig{bucket: bucket, prefix: prefix}
}

// DownloadTfConfigTo return bool apakah config di bucket ada atau ngga
// butuh tenant id dalam konteks
func (b *BucketConfig) DownloadTfConfigTo(ct context.Context, dir string) (bool, error) {

	client, err := storage.NewClient(ct)
	if err != nil {
		return false, fmt.Errorf("Failed to create client: %v\n", err)
	}
	defer client.Close()

	bucketName := b.bucket
	prefix := filepath.Join(b.prefix, ct.Value("tenant_id").(string))
	localDir := dir

	// Get the bucket handle
	bucket := client.Bucket(bucketName)

	// Create a query with the specified prefix
	query := &storage.Query{Prefix: prefix}

	// Initialize an iterator to list objects
	it := bucket.Objects(ct, query)
	is_exists := false
	// Iterate through the objects
	for {
		attrs, err := it.Next()
		if errors.Is(err, iterator.Done) {
			return is_exists, nil
		}
		if err != nil {
			return false, fmt.Errorf("Failed to list objects: %v\n", err)
		}

		// Download the object
		if err := downloadObject(ct, client, bucketName, attrs.Name, localDir); err != nil {
			return false, fmt.Errorf("Failed to download object %s: %v\n", attrs.Name, err)
		} else if strings.Contains(attrs.Name, "main.tf") {
			is_exists = true
		}
	}
}

func (b *BucketConfig) UploadTfConfig(ctx context.Context, dir string) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create upload client: %v\n", err)
	}
	bucketName := b.bucket
	remotePrefix := filepath.Join(b.prefix, ctx.Value("tenant_id").(string))
	localDir := dir

	defer client.Close()
	return filepath.WalkDir(localDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories themselves but not their contents
		if d.IsDir() {
			return nil
		}
		fmt.Printf("Walked to: %s\n", d.Name())
		// Determine the relative path of the file
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// Create the remote object name
		remotePath := filepath.Join(remotePrefix, relPath)
		remotePath = strings.ReplaceAll(remotePath, "\\", "/") // Ensure it's in Unix format

		// Upload the file
		if err := uploadFile(ctx, client, bucketName, path, remotePath); err != nil {
			return err
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
	localFilePath := filepath.Join(localDir, objectName)
	if err := os.MkdirAll(filepath.Dir(localFilePath), os.ModePerm); err != nil {
		return fmt.Errorf("mkdir: %v", err)
	}

	isDirectory := strings.HasSuffix(objectName, "/") // mungkin cuma perlu di cek belakangnya aja?
	if isDirectory {
		return nil // biar directory yang dibuat ga dijadiin file
	}

	localFile, err := os.Create(localFilePath)
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
