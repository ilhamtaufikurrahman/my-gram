package main

import (
	"my-gram/database"
	"my-gram/router"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(":8080")
}
