package log

type Format string

const (
	FormatText = "text"
	FormatJSON = "json"
)

type Config struct {
	Level      string        `yaml:"level"`
	Format     Format        `yaml:"format"`
	NoColor    bool          `yaml:"no_color"`
	ShowCaller bool          `yaml:"show_caller"`
	Sentry     *SentryConfig `yaml:"sentry,omitempty"`
}

type SentryConfig struct {
	DSN    string `yaml:"dsn"`
	Enable bool   `yaml:"enable"`
}
