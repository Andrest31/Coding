package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"st/coding" // импорт нашего пакета с кодированием
	_ "st/docs" // подключение сгенерированной swagger-документации

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	anotherRand "golang.org/x/exp/rand"
)

// Инициализируем глобальный источник случайных чисел
var randSrc anotherRand.Source

func init() {
	// Инициализация с текущим временем для обеспечения случайности
	randSrc = anotherRand.NewSource(uint64(time.Now().UnixMicro()))
}

// Структура входных и выходных данных
type DATA struct {
	Username       string `json:"username" binding:"required"`
	MessagePart    string `json:"message_part" binding:"required"`
	Timestamp      string `json:"timestamp" binding:"required"`
	SequenceNumber int    `json:"sequence_number" binding:"required"`
	TotalParts     int    `json:"total_parts" binding:"required"`
}

// Отправка данных обратно в транспортный уровень (предположительно другой сервис)
func SendCodeRequest(body DATA) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "main:8010/receive", bytes.NewBuffer(reqBody))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(reqBody)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка при отправке в транспортный уровень:", err)
		return
	}
	defer resp.Body.Close()
}

// @Summary Codes and decodes messages
// @Description Кодирует, вносит ошибки, исправляет и декодирует сообщение
// @Tags code
// @Accept json
// @Produce json
// @Param data body main.DATA true "Data to process"
// @Success 200 {object} main.DATA
// @Router /code [post]

func CodeHandler(c *gin.Context) {
	var data DATA

	if err := c.BindJSON(&data); err != nil {
		log.Println("Ошибка при разборе входных данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	processed, err := coding.ProcessMessage(data.MessagePart, randSrc)
	if err != nil {
		log.Println("Ошибка обработки сообщения:", err)
		c.JSON(http.StatusOK, gin.H{
			"error":        "Сообщение потеряно или повреждено",
			"sequence_num": data.SequenceNumber,
		})
		return
	}

	data.MessagePart = processed
	c.JSON(http.StatusOK, data)

	go SendCodeRequest(data)
}

func main() {
	router := gin.Default()

	// Эндпоинт для кодирования/декодирования сообщений
	router.POST("/code", CodeHandler)

	// Swagger-интерфейс по адресу /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Запуск HTTP-сервера на localhost:8081
	router.Run(":8020")
}
