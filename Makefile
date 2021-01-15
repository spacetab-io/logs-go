# ----
## LINTER stuff start
LINTER_VERSION=v1.27.0

get_lint_binary:
	@[ -f ./golangci-lint ] && echo "golangci-lint exists" || ( echo "getting golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./ $(LINTER_VERSION) && ./golangci-lint --version )
.PHONY: get_lint_binary

get_lint_config:
	@[ -f .golangci.yml ] && echo ".golangci.yml exists" || ( echo "getting .golangci.yml" && curl -O https://raw.githubusercontent.com/microparts/docker-golang/master/.golangci.yml )
.PHONY: get_lint_config

lint: get_lint_binary get_lint_config
	./golangci-lint run -v
.PHONY: lint

lint_quiet: get_lint_binary get_lint_config
	@./golangci-lint run
.PHONY: lint_quiet

## LINTER stuff end
# ----

# ----
## TEST stuff start

test-unit:
	go test $$(go list ./...) --race --cover -count=1 -timeout 1s -coverprofile=c.out -v
.PHONY: test-unit

coverage-html:
	go tool cover -html=c.out -o coverage.html
.PHONE: coverage-html

test: test-unit coverage-html
.PHONY: test

## TEST stuff end
# ----

circle:
	circleci local execute
