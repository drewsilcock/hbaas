package main

import (
	"log"

	"github.com/drewsilcock/hbaas-server/cmd"
)

// @title Happy Birthday as a Service
// @description This RESTful API says happy birthday to people by name and date.

// @host localhost:8000
// @BasePath /
// @schemes http

// @contact.name the Hartree Centre
// @contact.email drew.silcock@stfc.ac.uk

// @license.name MIT

// @tag.name Saying Happy Birthday
// @tag.description Endpoints for saying happy birthday to people given their name or person ID.

// @tag.name Person Management
// @tag.description Endpoints for CRUD operations on people.
func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal("Unable to run root command:", err)
	}
}
