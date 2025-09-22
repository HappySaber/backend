package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type message struct {
	message string `json:"message"`
}

func main() {
	port := "8080"
	router := gin.New()
	router.POST("/", addMessage)
	router.Run(":" + port)
}

func addMessage(c *gin.Context) {
	var mes message
	if err := c.ShouldBindJSON(&mes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("\n%s", mes.message))
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
	}
	c.IndentedJSON(http.StatusOK, nil)
}
