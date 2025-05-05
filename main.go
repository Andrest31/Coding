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
	Id            int    `json:"socket_id,omitempty"`               // ID сокета (опционально)
	Data          string `json:"data" binding:"required"`           // Полезная нагрузка
	SegmentNumber int    `json:"segment_number" binding:"required"` // Номер сегмента
	TotalSegments int    `json:"total_segments" binding:"required"` // Общее число сегментов
	Username      string `json:"username" binding:"required"`       // Имя пользователя
	SendTime      string `json:"send_time" binding:"required"`      // Время отправки
	MessageId     string `json:"message_id" binding:"required"`     // ID сообщения
}

// Отправка данных обратно в транспортный уровень (предположительно другой сервис)
func SendCodeRequest(body DATA) {
	reqBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "http://172.16.95.192:8081/transfer", bytes.NewBuffer(reqBody))
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
// @Schemes
// @Description Кодирует, вносит ошибки, исправляет и декодирует сообщение
// @Tags code
// @Param data body main.DATA true "data"
// @Accept json
// @Produce json
// @Success 200
// @Router /code [post]
func CodeHandler(c *gin.Context) {
	var data DATA

	// Пробуем распарсить JSON-запрос в структуру DATA
	if err := c.BindJSON(&data); err != nil {
		log.Println("Ошибка при разборе входных данных:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Обработка сообщения: кодирование, внесение ошибок, исправление
	processed, err := coding.ProcessMessage(data.Data, randSrc)
	if err != nil {
		log.Println("Ошибка обработки сообщения:", err)
		c.JSON(http.StatusOK, gin.H{
			"error":       "Сообщение потеряно или повреждено",
			"message_id":  data.MessageId,
			"segment_num": data.SegmentNumber,
		})
		return
	}

	// Возвращаем обработанное сообщение обратно
	data.Data = processed
	c.JSON(http.StatusOK, data)

	// Отправляем сообщение на транспортный уровень в фоне
	go SendCodeRequest(data)
}

func main() {
	router := gin.Default()

	// Эндпоинт для кодирования/декодирования сообщений
	router.POST("/code", CodeHandler)

	// Swagger-интерфейс по адресу /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Запуск HTTP-сервера на localhost:8081
	router.Run("127.0.0.1:8081")
}
