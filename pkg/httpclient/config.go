package httpclient

type Config struct {
	BaseURL     string `mapstructure:"base_url" yaml:"base_url"`
	EnableTrace bool   `mapstructure:"enable_trace" yaml:"enable_trace"`
}
