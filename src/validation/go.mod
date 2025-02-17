module github.com/TibiaData/tibiadata-api-go/src/validation

go 1.24.0

replace github.com/TibiaData/tibiadata-api-go/src/tibiamapping => ../tibiamapping

require (
	github.com/TibiaData/tibiadata-api-go/src/tibiamapping v0.0.0-20250214153538-6bdfa7545f05
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
