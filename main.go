package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"errors"
	"github.com/containrrr/shoutrrr"
	t "github.com/containrrr/shoutrrr/pkg/types"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slices"
)

type Notification struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func getToken(c *gin.Context) (string, error) {
	// TODO: error handling
	reqToken := c.GetHeader("Authorization")
	if len(reqToken) == 0 {
		return "", errors.New("No key provided in Authorization header")
	}

	splitToken := strings.Split(reqToken, "Bearer ")
	return splitToken[1], nil
}

func verifyToken(c *gin.Context) (bool, error) {
	// TODO: error handling
	token, err := getToken(c)
	if err != nil {
		return false, err
	}

	keys := strings.Split(os.Getenv("KEYS"), ",")

	return slices.Contains(keys, token), nil
}

func postNotification(c *gin.Context) {
	valid, err := verifyToken(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	if !valid {
		c.String(http.StatusUnauthorized, "Invalid key in Authorization header")
		return
	}

	var notification Notification

	if err := c.BindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := fmt.Sprintf("discord://%s@%s", os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_WEBHOOK_ID"))
	sender, err := shoutrrr.CreateSender(url)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating sender")
		return
	}

	var params t.Params = map[string]string{"title": notification.Title}
	errs := sender.Send(notification.Message, &params)

	if len(errs) != 0 && errs[0] != nil {
		c.String(http.StatusInternalServerError, "Error sending notification")
		return
	}

	// Add the new album to the slice.
	c.String(http.StatusOK, "Successfully sent notification")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("[WARN] Failed to load .env file")
	}

	router := gin.Default()
	router.POST("/send", postNotification)

	router.Run("0.0.0.0:8080")
}
