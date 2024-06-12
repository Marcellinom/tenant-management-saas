package terraformProduct

import (
	"github.com/go-git/go-git"
	"os"
)

type ProductBackend interface {
	CloneTo(path string) error
	DeleteOn(path string) error
	GetProductConfig() *ProductConfig
}

type ProductStoredOnGit struct {
	product_config *ProductConfig
}

func UsingGit(config *ProductConfig) ProductStoredOnGit {
	return ProductStoredOnGit{product_config: config}
}

func (g ProductStoredOnGit) CloneTo(path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      g.product_config.product_url,
		Progress: os.Stdout,
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: implement this
func (g ProductStoredOnGit) DeleteOn(path string) error {
	return nil
}

func (g ProductStoredOnGit) GetProductConfig() *ProductConfig {
	return g.product_config
}
