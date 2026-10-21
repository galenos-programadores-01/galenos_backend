// Package httpadapter es el adaptador de entrada (driving adapter) que
// expone los casos de uso del dominio como una REST API para el frontend
// Angular. No contiene reglas de negocio: solo traduce HTTP <-> puertos.
package httpadapter

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// apiResponse es el sobre estándar de toda respuesta de la API, para que
// Angular pueda manejar éxito y error de forma consistente.
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
	mensajeCliente := message
	if status >= http.StatusInternalServerError {
		log.Printf("ERROR SERVIDOR [%s]: %s", code, message)
		mensajeMinuscula := strings.ToLower(message)
		if strings.Contains(mensajeMinuscula, "mssql:") || strings.Contains(mensajeMinuscula, "sql:") {
			mensajeCliente = "Ocurrio un error en el procesamiento de la base de datos"
		}
	}
	c.JSON(status, apiResponse{Success: false, Error: &apiError{Code: code, Message: mensajeCliente}})
}
