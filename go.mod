module github.com/TibiaData/tibiadata-api-go

go 1.23

replace github.com/TibiaData/tibiadata-api-go/src/tibiamapping => ./src/tibiamapping

replace github.com/TibiaData/tibiadata-api-go/src/static => ./src/static

replace github.com/TibiaData/tibiadata-api-go/src/validation => ./src/validation

require (
	github.com/PuerkitoBio/goquery v1.10.0
	github.com/TibiaData/tibiadata-api-go/src/static v0.0.0-20240821085055-c65129c5697a
	github.com/TibiaData/tibiadata-api-go/src/validation v0.0.0-20240821085055-c65129c5697a
	github.com/gin-contrib/gzip v1.0.1
	github.com/gin-gonic/gin v1.10.0
	github.com/go-resty/resty/v2 v2.15.3
	github.com/mantyr/go-charset v0.0.0-20160510214718-44d054d82c4a
	github.com/stretchr/testify v1.9.0
	golang.org/x/text v0.20.0
)

require (
	github.com/TibiaData/tibiadata-api-go/src/tibiamapping v0.0.0-20240821085055-c65129c5697a // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/bytedance/sonic v1.12.1 // indirect
	github.com/bytedance/sonic/loader v0.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.5 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	golang.org/x/arch v0.9.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
