logs-go
-------

[![CircleCI](https://circleci.com/gh/spacetab-io/logs-go.svg?style=shield)](https://circleci.com/gh/spacetab-io/logs-go) [![codecov](https://codecov.io/gh/spacetab-io/logs-go/graph/badge.svg)](https://codecov.io/gh/spacetab-io/logs-go)

Wrapper for [uber zap](go.uber.org/zap) logger tuned to work with [configuration](https://github.com/spacetab-io/configuration-go) and
sentry hook.

## Usage

Initiate new logger filled with struct that implements `LogsConfigInterface` and use it as common zap logger

```go
package main

import (
	"os"

	cfgstructs "github.com/spacetab-io/configuration-structs-go/v2"
	log "github.com/spacetab-io/logs-go/v3"
)

func main() {
	conf := &cfgstructs.Logs{
		Level:   "debug",
		Format:  "text",
		Colored: true,
		Caller:  cfgstructs.CallerConfig{Show: true, SkipFrames: 1},
		Sentry: &cfgstructs.SentryConfig{
			Enable: true,
			Debug:  true,
			DSN:    os.Getenv("SENTRY_DSN"),
		},
	}

	logger, err := log.Init(conf, "test", "service", "v3.0.0", os.Stdout)
	if err != nil {
		panic(err)
	}

	logger.Warn().Msg("log some warning")
}
```

## Licence

The software is provided under [MIT Licence](LICENCE).