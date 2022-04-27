module github.com/spacetab-io/logs-go/v3

go 1.17

require (
	github.com/getsentry/sentry-go v0.13.0
	github.com/json-iterator/go v1.1.12
	github.com/spacetab-io/configuration-structs-go v0.1.1
	github.com/stretchr/testify v1.7.1
	go.uber.org/zap v1.21.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/spacetab-io/configuration-structs-go => ../configuration-structs-go
