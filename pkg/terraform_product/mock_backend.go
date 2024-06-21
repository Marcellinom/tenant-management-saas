package terraform_product

type MockProductBackend struct {
}

func Mock() MockProductBackend {
	return MockProductBackend{}
}

func (p MockProductBackend) CopyTo(path string) error {
	return nil
}

func (p MockProductBackend) GetProductConfig() *ProductConfig {
	return &ProductConfig{}
}
