package api

import (
	"./controllers"
	"./seed"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var server = controllers.Server{}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("MYSQL_URL"))
	seed.Load(server.DB)
	server.Run(os.Getenv("PORT"))

}
