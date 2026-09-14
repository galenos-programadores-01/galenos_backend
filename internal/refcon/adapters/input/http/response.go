// Package refconhttp es el adaptador de entrada (REST) del módulo de
// referencias. Usa el mismo sobre estándar de respuesta de la API
// (success/data/error) pero de forma autocontenida para no acoplarse a la
// capa HTTP base.
package refconhttp

import "github.com/gin-gonic/gin"

// apiResponse es el sobre estándar de toda respuesta de la API.
type apiResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *apiError `json:"error,omitempty"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func respondSuccess(c *gin.Context, status int, data any) {
	c.JSON(status, apiResponse{Success: true, Data: data})
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, apiResponse{Success: false, Error: &apiError{Code: code, Message: message}})
}
