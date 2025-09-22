package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Message string `json:"message"`
}

func main() {
	port := "8080"
	router := gin.New()

	router.Use(corsMiddleware())

	router.POST("/", addMessage)
	router.GET("/", getMessages)
	router.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func getMessages(c *gin.Context) {
	data, err := os.ReadFile("data.txt")
	if err != nil {
		if os.IsNotExist(err) {
			c.String(http.StatusOK, "")
			return
		}
		fmt.Println("Ошибка чтения файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}
	c.String(http.StatusOK, string(data))
}

func addMessage(c *gin.Context) {
	var mes Message
	if err := c.ShouldBindJSON(&mes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Получено сообщение: %s\n", mes.Message)

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File error"})
		return
	}
	defer file.Close()

	_, err = file.WriteString(mes.Message + "\n")
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Write error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Data saved"})
}
