package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type Server struct {
	DB     *gorm.DB
	Router *gin.Engine
}

func (server *Server) Initialize(DBDriver, DBUrl string) {
	var err error
	if DBDriver == "mysql" {
		server.DB, err = gorm.Open(DBDriver, DBUrl)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Println("We are connected to the", DBDriver, "database")
		}
	}
	if DBDriver == "postgres" {
		server.DB, err = gorm.Open(DBDriver, DBUrl)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", DBDriver)
			fmt.Println()
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", DBDriver)
			fmt.Println()
		}
	}
	server.Router = gin.Default()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to localhost" + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
