package main

import (
	"gcp-access-token/initial"
	"gcp-access-token/routes"
)

func main() {
	initial.InitEnvConfigs()
	routes.SetupRouter()
}
