logs-go
-------

[![CircleCI](https://circleci.com/gh/spacetab-io/logs-go.svg?style=shield)](https://circleci.com/gh/spacetab-io/logs-go) [![codecov](https://codecov.io/gh/spacetab-io/logs-go/graph/badge.svg)](https://codecov.io/gh/spacetab-io/logs-go)

[Logrus](github.com/sirupsen/logrus) wrapper for easy use with sentry hook, database (gorm) and mux (gin) loggers.

## Usage

Initiate new logger with filled `logs.Config` and use it as common logrus logger instance

```go
package main

import (
	"time"
	
	"github.com/spacetab-io/logs-go"
)

func main() {
	conf := &logs.Config{
		LogLevel: "warn",
		Debug: true,
		Sentry: &logs.SentryConfig{
			Enable: true,
			DSN: "http://dsn.sentry.com",
		},
	}
	
	l, err := logs.NewLogger(conf)
	if err != nil {
		panic(err)
	}
	
	l.Warn("log some warning")
}
```

## Licence

The software is provided under [MIT Licence](LICENCE).