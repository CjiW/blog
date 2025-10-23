- golangci-lint
```yml
# curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin
version: "2"
linters:
  enable:
    - asasalint
    - bidichk
    - bodyclose
    - containedctx
    - gocheckcompilerdirectives
    - makezero
    - misspell
    - nilerr
    - nolintlint
    - nosprintfhostport
    - unconvert
    - usetesting
    - wastedassign
    - whitespace
  disable:
    - errcheck
    - usestdlibvars
  settings:
    staticcheck:
      checks:
        - all
        - -SA1019
  exclusions:
    generated: lax
    presets:
      - comments
      - common-false-positives
      - legacy
      - std-error-handling
    paths:
      - third_party$
      - builtin$
      - examples$
severity:
  default: error
  rules:
    - linters:
        - gofmt
        - goimports
        - intrange
      severity: info
formatters:
  enable:
    - gofmt
    - gofumpt
  exclusions:
    generated: lax
    paths:
      - third_party$
      - builtin$
      - examples$
```
