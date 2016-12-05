package main

import (

	"os"

	"github.com/pineda89/golang-springboot/eureka"
	"github.com/pineda89/golang-springboot/actuator"
	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/config"
	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/service"
	"github.com/oneboxtm/integrations-platform/int-avet-conciliation/handler"
)



func main() {

	// Load configuration in order to start application
	config.LoadConfig()

	// Controller to handle application webservice endpoints
	go handler.StartWebServer( config.Configuration["server.port"].(int) )

	// Start actuator webservices endpoints (/info, /health)
	go actuator.InitializeActuator()

	// Start application services
	go service.StartService()

	// Register to the service discovery
	go eureka.Register(config.Configuration)
	eureka.CaptureInterruptSignal()
	eureka.Deregister()
	os.Exit(0)

}
