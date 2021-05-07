package main

import (

)

var (
	// Drone Server Information
	DroneServers = []DroneServer {
		DroneServer {
			Endpoint: "http://10.10.25.16:9000",
			YamlPluginSecret: "bea26a2221fd8090ea38720fc445eca6",
			WebhookPluginSecret: "bea26a2221fd8090ea38720fc445eca6",
		},
		DroneServer {
			Endpoint: "http://10.16.2.37:9000",
			YamlPluginSecret: "bea26a2221fd8090ea38720fc445eca6",
			WebhookPluginSecret: "bea26a2221fd8090ea38720fc445eca6",
		},
	}

	// Rolling Information
	RollingEndpoint = "http://10.16.2.37:8080"
)

var (
	Rolling = NewRollingClient(RollingEndpoint)
)