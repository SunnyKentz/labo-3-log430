package utils

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/models"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FAILURE_ERR(err error) string { return "Failed Request : " + err.Error() }
func SYNTAX_ERR(source string, value any) string {
	return fmt.Sprintf("%s : %v could not be parsed, check syntax", source, value)
}
func NOTFOUND_ERR(source string, value any) string {
	return fmt.Sprintf("%s : %v was not found", source, value)
}

func GetApiError(c *fiber.Ctx, message string, status int) error {
	errorMsg := ""
	switch status {
	case 400:
		errorMsg = "Bad Request"
	case 401:
		errorMsg = "Unauthorized"
	case 403:
		errorMsg = "Forbidden"
	case 404:
		errorMsg = "Not Found "
	case 500:
		errorMsg = "Internal Server Error"
	}
	return c.JSON(models.ApiError{
		Timestamp: time.Now(),
		Status:    status,
		Error:     errorMsg,
		Message:   message,
		Path:      c.Route().Path,
	})
}

func GetApiSuccess(c *fiber.Ctx, status int) error {
	message := ""
	switch status {
	case 200:
		message = "Ok"
	case 201:
		message = "Created"
	}
	return c.JSON(models.ApiSuccess{
		Timestamp: time.Now(),
		Status:    status,
		Message:   message,
		Path:      c.Route().Path,
	})
}

func API_MAGASIN() string  { return "http://" + os.Getenv("GATEWAY") + ":8080" }
func API_MERE() string     { return "http://" + os.Getenv("GATEWAY") + ":8090" }
func API_LOGISTIC() string { return "http://" + os.Getenv("GATEWAY") + ":8091" }

func Errnotnil(err error) {
	if err != nil {
		logger.Error(err.Error())
	}
}
