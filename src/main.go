package main

import "log"

var (
	// TibiaDataDefaultVoc - default vocation when not specified in request
	TibiaDataDefaultVoc string = "all"

	// TibiaData app flags for running
	TibiaDataAPIversion int = 3
	TibiaDataDebug      bool

	// TibiaData app settings
	TibiaDataHost string // set through env TIBIADATA_HOST

	// TibiaData app details set to release/build on GitHub
	TibiaDataBuildRelease = "unknown"     // will be set by GitHub Actions (to release number)
	TibiaDataBuildBuilder = "manual"      // will be set by GitHub Actions
	TibiaDataBuildCommit  = "-"           // will be set by GitHub Actions (to git commit)
	TibiaDataBuildEdition = "open-source" //
)

// @title           TibiaData API
// @version         edge
// @description     This is the API documentation for the TibiaData API.
// @description     The documentation contains version 3 and above.
// @termsOfService  https://tibiadata.com/terms/

// @contact.name   TibiaData
// @contact.url    https://www.tibiadata.com/contact/
// @contact.email  tobias@tibiadata.com

// @license.name  MIT
// @license.url   https://github.com/TibiaData/tibiadata-api-go/blob/main/LICENSE

// @schemes   http
// @host      localhost:8080
// @BasePath  /

func main() {
	// logging start of TibiaData
	log.Printf("[info] TibiaData API starting..")

	// running the TibiaDataInitializer function
	TibiaDataInitializer()

	// logging build information
	log.Printf("[info] TibiaData API release: %s", TibiaDataBuildRelease)
	log.Printf("[info] TibiaData API build: %s", TibiaDataBuildBuilder)
	log.Printf("[info] TibiaData API commit: %s", TibiaDataBuildCommit)
	log.Printf("[info] TibiaData API edition: %s", TibiaDataBuildEdition)

	// setting tibiadata-application to log much less if DEBUG_MODE is false (default is false)
	if !getEnvAsBool("DEBUG_MODE", false) {
		log.Printf("[info] TibiaData API debug-mode: disabled")
	} else {
		// setting debug to true for more logging
		TibiaDataDebug = true
		log.Printf("[info] TibiaData API debug-mode: enabled")

		// logging user-agent string
		log.Printf("[debug] TIbiaData API User-Agent: %s", TibiaDataUserAgent)
	}

	// starting webserver.go stuff
	runWebServer()
}

// TibiaDataInitializer func - init things at beginning
func TibiaDataInitializer() {
	// setting TibiaDataBuildEdition
	if isEnvExist("TIBIADATA_EDITION") {
		TibiaDataBuildEdition = getEnv("TIBIADATA_EDITION", "open-source")
	}

	// adding information of host
	TibiaDataHost = getEnv("TIBIADATA_HOST", "")
	if TibiaDataHost != "" {
		TibiaDataHost = "+https://" + TibiaDataHost
	}

	// generating TibiaDataUserAgent with TibiaDataUserAgentGenerator function
	TibiaDataUserAgent = TibiaDataUserAgentGenerator(TibiaDataAPIversion)

	// setting TibiaDataProxyDomain
	if isEnvExist("TIBIADATA_PROXY") {
		TibiaDataProxyDomain = "https://" + getEnv("TIBIADATA_PROXY", "www.tibia.com") + "/"
	}

	log.Printf("[info] TibiaData API proxy: %s", TibiaDataProxyDomain)

	// initializing houses mappings
	TibiaDataHousesMappingInitiator()

	// run some functions that are empty but required for documentation to be done
	_ = tibiaNewslistArchiveV3()
	_ = tibiaNewslistArchiveDaysV3()
	_ = tibiaNewslistLatestV3()
}
