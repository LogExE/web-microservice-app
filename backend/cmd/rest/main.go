package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/LogExE/web-microservice-app/config"
	"github.com/LogExE/web-microservice-app/internal"
)

type app struct {
	boxRepo boxes.BoxRepo
}

func main() {
	cfg := config.New()

	db, err := boxes.InitDB()
	if err != nil {
		log.Fatal("Failed to init db: ", err)
	}
	p, err := boxes.MakeProducer(cfg, "boxesMAIN")
	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}

	a := app{
		boxRepo: boxes.BoxRepo{DB: db},
	}

	go outboxWorker(p, &a.boxRepo)

	router := gin.Default()
	a.registerRoutes(router)
	router.Run(cfg.REST.Addr)
}

func (a *app) registerRoutes(r *gin.Engine) {
	r.GET("/boxes", a.boxesGet)

	r.GET("/box/:id", a.boxGet)
	r.POST("/box", a.boxPost)

	r.POST("/like/:id", a.boxLike)
}
