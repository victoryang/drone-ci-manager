package main

import (

)

const (
	// Drone Plugin Information
	YamlPluginSecret = "bea26a2221fd8090ea38720fc445eca6"
	WebhookPluginSecret = "bea26a2221fd8090ea38720fc445eca6"

	// Drone Server Information
	DroneServerOwner = "Test"
	DroneServerEndpoint = "http://10.10.25.16:9000"
	DroneServerToken = "9sZnUndOShuMAj9zViZNdwaSFl4ovkTj"

	// Rolling Information
	RollingEndpoint = "http://10.10.25.16:8080"
)

var (
	DroneCli = NewDroneClient(DroneServerOwner, DroneServerHost, DroneServerToken)	
	Rolling = NewRollingClient(RollingEndpoint)
)