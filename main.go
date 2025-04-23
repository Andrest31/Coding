package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"st/coding"
	"strconv"
	"time"

	_ "st/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	anotherRand "golang.org/x/exp/rand"
)

var randSrc anotherRand.Source

func init() {
	randSrc = anotherRand.NewSource(uint64(time.Now().UnixMicro()))
}

type DATA struct {
	Id            int    `json:"socket_id,omitempty"`
	Data          string `json:"data" binding:"required"`
	SegmentNumber int    `json:"segment_number" binding:"required"`
	TotalSegments int    `json:"total_segments" binding:"required"`
	Username      string `json:"username" binding:"required"`
	SendTime      string `json:"send_time" binding:"required"`
	MessageId     string `json:"message_id" binding:"required"`
}

func SendCodeRequest(body DATA) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "http://172.16.95.192:8080/transfer", bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}

// Code godoc
// @Summary Codes and decodes messages
// @Schemes  main.DATA
// @Description  Codes and decodes messages
// @Tags code
// @Param		data body main.DATA true  "data"
// @Accept json
// @Produce json
// @Success 200
// @Router /code [post]
func CodeHandler(c *gin.Context) {

	var data DATA
	err := c.BindJSON(&data)
	if err != nil {
		log.Println("ERROR__", err)
		fmt.Println(time.Now())
	}
	msg, err := coding.ProcessMessage(data.Data, randSrc)
	if err != nil {
		fmt.Println("Сообщение утеряно")

	} else {
		data.Data = msg
		c.JSON(http.StatusOK, data)
		go SendCodeRequest(data)
	}

}

// @title           Swagger API
// @version         1.0
// @description    This is a coding server for Networking.

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	router := gin.Default()

	router.POST("/code", CodeHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.Run("127.0.0.1:8080")

}
