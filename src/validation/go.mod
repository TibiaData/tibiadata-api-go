module github.com/TibiaData/tibiadata-api-go/src/validation

go 1.23

replace github.com/TibiaData/tibiadata-api-go/src/tibiamapping => ../tibiamapping

require (
	github.com/TibiaData/tibiadata-api-go/src/tibiamapping v0.0.0-20240821085055-c65129c5697a
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-resty/resty/v2 v2.14.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
