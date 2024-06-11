package terraformProduct

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

// TODO: implement this
func (g ProductStoredOnGit) CloneTo(path string) error {
	return nil
}

// TODO: implement this
func (g ProductStoredOnGit) DeleteOn(path string) error {
	return nil
}

func (g ProductStoredOnGit) GetProductConfig() *ProductConfig {
	return g.product_config
}
