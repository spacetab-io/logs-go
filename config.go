package log

type Format string

const (
	FormatText = "text"
	FormatJSON = "json"
)

type CallerConfig struct {
	Disabled         bool `yaml:"hide_caller"`
	CallerSkipFrames int  `yaml:"skip_frames"`
}

type Config struct {
	Level   string        `yaml:"level"`
	Format  Format        `yaml:"format"`
	NoColor bool          `yaml:"no_color"`
	Caller  *CallerConfig `yaml:"caller"`
	Sentry  *SentryConfig `yaml:"sentry,omitempty"`
}

type SentryConfig struct {
	DSN    string `yaml:"dsn"`
	Enable bool   `yaml:"enable"`
}
