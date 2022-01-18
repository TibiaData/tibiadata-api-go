package main

import "log"

var (
	// TibiadataDefaultVoc - default vocation when not specified in request
	TibiadataDefaultVoc string = "all"

	// Tibiadata app flags for running
	TibiadataAPIversion int = 3
	TibiadataDebug      bool

	// Tibiadata app settings
	TibiadataHost string // set through env TIBIADATA_HOST

	// Tibiadata app details set to release/build on GitHub
	TibiadataBuildRelease = "unknown"     // will be set by GitHub Actions (to release number)
	TibiadataBuildBuilder = "manual"      // will be set by GitHub Actions
	TibiadataBuildCommit  = "-"           // will be set by GitHub Actions (to git commit)
	TibiadataBuildEdition = "open-source" //
)

// @title           TibiaData API
// @version         edge
// @description     This is the API documentation for the TibiaData API written in Golang.
// @termsOfService  https://tibiadata.com/terms/

// @contact.name   TibiaData Support
// @contact.url    https://tibiadata.com/contact/
// @contact.email  tobias@tibiadata.com

// @license.name  MIT
// @license.url   https://github.com/TibiaData/tibiadata-api-go/blob/main/LICENSE

// @host      dev.tibiadata.com
// @BasePath  /

func main() {
	// logging start of TibiaData
	log.Printf("[info] TibiaData API starting..")

	// running the TibiaDataInitializer function
	TibiaDataInitializer()

	// logging build information
	log.Printf("[info] TibiaData API release: %s", TibiadataBuildRelease)
	log.Printf("[info] TibiaData API build: %s", TibiadataBuildBuilder)
	log.Printf("[info] TibiaData API commit: %s", TibiadataBuildCommit)
	log.Printf("[info] TibiaData API edition: %s", TibiadataBuildEdition)

	// setting tibiadata-application to log much less if DEBUG_MODE is false (default is false)
	if !getEnvAsBool("DEBUG_MODE", false) {
		log.Printf("[info] TibiaData API debug-mode: disabled")
	} else {
		// setting debug to true for more logging
		TibiadataDebug = true
		log.Printf("[info] TibiaData API debug-mode: enabled")

		// logging user-agent string
		log.Printf("[debug] TIbiaData API User-Agent: %s", TibiadataUserAgent)
	}

	// starting webserver.go stuff
	runWebServer()
}

// TibiaDataInitializer func - init things at beginning
func TibiaDataInitializer() {
	// setting TibiadataBuildEdition
	if isEnvExist("TIBIADATA_EDITION") {
		TibiadataBuildEdition = getEnv("TIBIADATA_EDITION", "open-source")
	}

	// adding information of host
	TibiadataHost = getEnv("TIBIADATA_HOST", "")
	if TibiadataHost != "" {
		TibiadataHost = "+https://" + TibiadataHost
	}

	// generating TibiadataUserAgent with TibiadataUserAgentGenerator function
	TibiadataUserAgent = TibiadataUserAgentGenerator(TibiadataAPIversion)

	// setting TibiadataProxyDomain
	if isEnvExist("TIBIADATA_PROXY") {
		TibiadataProxyDomain = "https://" + getEnv("TIBIADATA_PROXY", "www.tibia.com") + "/"
	}

	log.Printf("[info] TibiaData API proxy: %s", TibiadataProxyDomain)

	// initializing houses mappings
	TibiaDataHousesMappingInitiator()

	// run some functions that are empty but required for documentation to be done
	tibiaNewslistArchiveV3()
	tibiaNewslistArchiveDaysV3()
	tibiaNewslistLatestV3()
}
