package terraformProduct

type ProductConfig struct {
	product_url           string
	tf_entrypoint_url     string
	script_entrypoint_url string
}

func (p ProductConfig) GetProductUrl() string {
	return p.product_url
}

func (p ProductConfig) GetTfEntrypoint() string {
	return p.tf_entrypoint_url
}

func (p ProductConfig) GetScriptEntrypoint() string {
	return p.script_entrypoint_url
}

func NewProductConfig(product_url, tf_entrypoint_url, script_entrypoint_url string) *ProductConfig {
	return &ProductConfig{
		product_url:           product_url,
		tf_entrypoint_url:     tf_entrypoint_url,
		script_entrypoint_url: script_entrypoint_url,
	}
}
