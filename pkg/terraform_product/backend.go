package terraform_product

import (
	"gopkg.in/src-d/go-git.v4"
	"os"
	"path/filepath"
	"strings"
)

type ProductBackend interface {
	CopyTo(path string) error
	GetProductConfig() *ProductConfig
}

type ProductStoredOnGit struct {
	product_config *ProductConfig
}

func UsingGit(config *ProductConfig) ProductStoredOnGit {
	return ProductStoredOnGit{product_config: config}
}

func (g ProductStoredOnGit) CopyTo(path string) error {
	_, err := git.PlainClone(filepath.Join(path, g.repoName(g.product_config.product_url)), false, &git.CloneOptions{
		URL:      g.product_config.product_url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g ProductStoredOnGit) GetProductConfig() *ProductConfig {
	return g.product_config
}

func (g ProductStoredOnGit) repoName(path string) string {
	without_git_prefix := strings.Replace(path, ".git", "", -1)
	lastSlash := strings.LastIndexAny(without_git_prefix, string(filepath.Separator)+"/")

	if lastSlash == -1 {
		// If no slash is found, return the entire string
		return without_git_prefix
	}

	// Return the substring from the last slash to the end of the string
	return without_git_prefix[lastSlash+1:]
}
