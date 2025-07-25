module github.com/eiiches/mysql-protobuf-functions

go 1.24.4

require (
	github.com/bufbuild/protocompile v0.14.1
	github.com/go-sql-driver/mysql v1.9.3
	github.com/google/go-github/v66 v66.0.0
	github.com/onsi/gomega v1.37.0
	github.com/samber/lo v1.51.0
	github.com/urfave/cli/v3 v3.3.8
	google.golang.org/protobuf v1.36.5
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/mna/pigeon v1.3.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	mvdan.cc/gofumpt v0.8.0 // indirect
)

tool (
	github.com/mna/pigeon
	mvdan.cc/gofumpt
)
