package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/LogExE/web-microservice-app/internal"
)

type app struct {
	boxRepo boxes.BoxRepo
}

func main() {
	router := gin.Default()

	db, err := boxes.InitDB()
	if err != nil {
		log.Fatal("Failed to init db: ", err)
	}
	p, err := boxes.MakeProducer("boxesMAIN")
	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}

	a := app{
		boxRepo: boxes.BoxRepo{DB: db},
	}

	go outboxWorker(p, &a.boxRepo)

	a.registerRoutes(router)
	router.Run("localhost:8080")
}

func (a *app) registerRoutes(r *gin.Engine) {
	r.GET("/boxes", a.boxesGet)

	r.GET("/box/:id", a.boxGet)
	r.POST("/box", a.boxPost)

	r.POST("/like/:id", a.boxLike)
}
