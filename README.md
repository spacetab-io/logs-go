logs-go
-------

[![CircleCI](https://circleci.com/gh/spacetab-io/logs-go.svg?style=shield)](https://circleci.com/gh/spacetab-io/logs-go) [![codecov](https://codecov.io/gh/spacetab-io/logs-go/graph/badge.svg)](https://codecov.io/gh/spacetab-io/logs-go)

Wrapper for [zerolog](https://github.com/rs/zerolog) tuned to work with [configuration](https://github.com/spacetab-io/configuration-go) and sentry hook.

## Usage

Initiate new logger with filled `log.Config` and use it as common zerolog

```go
package main

import (
	"github.com/spacetab-io/logs-go/v2"
)

func main() {
	conf := log.Config{
		Level:      "warn",
		Format:     "text",
		ShowCaller: true,
		Sentry: &log.SentryConfig{
			Enable: true,
			DSN:    "http://dsn.sentry.com",
		},
	}

	if err := log.Init("test", conf, "logs-go", "v2.*.*", nil); err != nil {
		panic(err)
	}

	log.Warn().Msg("log some warning")
}
```

## Licence

The software is provided under [MIT Licence](LICENCE).