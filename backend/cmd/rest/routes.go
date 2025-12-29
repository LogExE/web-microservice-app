package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/LogExE/web-microservice-app/internal"
)

func (a *app) boxesGet(c *gin.Context) {
	sortBy := c.Query("sortBy")
	
	boxes, err := a.boxRepo.BoxesList()
	if err != nil {
		log.Println("boxesGet ISE: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ISE"})
	}
	if sortBy == "likes" {
		sort.Slice(boxes, func(i, j int) bool {
			return boxes[i].Likes > boxes[j].Likes
		})
	}
	c.JSON(http.StatusOK, boxes)
}

func (a *app) boxPost(c *gin.Context) {
	var newBox boxes.Box
	if err := c.BindJSON(&newBox); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "json error"})
		return
	}

	newBox.Likes = 0
	err := a.boxRepo.BoxInsert(&newBox)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "insert error"})
		return
	}
	c.JSON(http.StatusCreated, newBox)
}

func (a *app) boxGet(c *gin.Context) {
	id, succ := intParamOrBadRequest("id", c)
	if !succ {
		return
	}

	box, err := a.boxRepo.BoxRetrieve(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "box doesn't exist"})
			return
		}
		log.Println("boxGet ISE: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ISE"})
		return
	}

	c.JSON(http.StatusOK, box)
}

func (a *app) boxLike(c *gin.Context) {
	id, succ := intParamOrBadRequest("id", c)
	if !succ {
		return
	}

	err := a.boxRepo.BoxLikeUpdate(int64(id))
	if err != nil {
		log.Println("boxLike ISE: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ISE"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func intParamOrBadRequest(param string, c *gin.Context) (int, bool) {
	p := c.Param(param)
	conv, err := strconv.Atoi(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"mesage": "bad int param: " + param})
		return 0, false
	}
	return conv, true
}
