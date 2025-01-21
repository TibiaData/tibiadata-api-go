module github.com/TibiaData/tibiadata-api-go/src/validation

go 1.23.5

replace github.com/TibiaData/tibiadata-api-go/src/tibiamapping => ../tibiamapping

require (
	github.com/TibiaData/tibiadata-api-go/src/tibiamapping v0.0.0-20250121083616-f3d6d188fdf7
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-resty/resty/v2 v2.16.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
