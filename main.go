package main

import (
	"challange10/database"
	"challange10/router"
)

func main() {
	database.StartDB()
	r := router.StartApp()
	r.Run(":8080")
}
