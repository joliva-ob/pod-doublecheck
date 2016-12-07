package main

import (

	"os"

	"github.com/pineda89/golang-springboot/eureka"
	"github.com/pineda89/golang-springboot/actuator"
	"github.com/joliva-ob/pod-doublecheck/config"
	"github.com/joliva-ob/pod-doublecheck/service"
	"github.com/joliva-ob/pod-doublecheck/handler"
)



func main() {

	// Load configuration in order to start application
	config.LoadConfig()

	// Controller to handle application webservice endpoints (/metrics)
	go handler.StartWebServer( config.Configuration["server.port"].(int) )

	// Start actuator webservices endpoints (/info, /health)
	go actuator.InitializeActuator()

	// Start application services
	go service.StartService()

	// Register to the service discovery
	go eureka.Register(config.Configuration)

	config.Log.Noticef("Application successfully started.")
	eureka.CaptureInterruptSignal()
	eureka.Deregister()
	os.Exit(0)

}
