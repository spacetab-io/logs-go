logs-go
-------

[![CircleCI](https://circleci.com/gh/spacetab-io/logs-go.svg?style=shield)](https://circleci.com/gh/spacetab-io/logs-go) [![codecov](https://codecov.io/gh/spacetab-io/logs-go/graph/badge.svg)](https://codecov.io/gh/spacetab-io/logs-go)

Wrapper for [zerolog](https://github.com/rs/zerolog) tuned to work with [configuration](https://github.com/spacetab-io/configuration-go) and
sentry hook.

## Usage

Initiate new logger filled with struct that implements `LogsConfigInterface` and use it as common zerolog

```go
package main

import (
	"os"

	cfgstructs "github.com/spacetab-io/configuration-structs-go"
	log "github.com/spacetab-io/logs-go/v3"
)

func main() {
	conf := cfgstructs.Logs{
		Level:  "warn",
		Format: cfgstructs.LogFormatText,
		Caller: &cfgstructs.CallerConfig{
			Disabled:         false,
			CallerSkipFrames: 2,
		},
		Sentry: &cfgstructs.SentryConfig{
			Enable: true,
			DSN:    "http://dsn.sentry.com",
		},
	}

	logger, err := log.Init("test", conf, "serviceName", "v3.1.2", os.Stdout)
	if err != nil {
		panic(err)
	}

	logger.Warn().Msg("log some warning")
}
```

## Licence

The software is provided under [MIT Licence](LICENCE).