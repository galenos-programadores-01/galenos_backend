package httpadapter

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// AuthHandler expone el puerto de entrada input.AuthService.
type AuthHandler struct {
	service input.AuthService
}

// NewAuthHandler inyecta el caso de uso en el adaptador HTTP.
func NewAuthHandler(service input.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login maneja POST /api/v1/auth/login y retorna un JWT Bearer que debe
// enviarse en el header "Authorization: Bearer <token>" del resto de
// peticiones.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_REQUEST", "username y password son requeridos")
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			respondError(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", err.Error())
			return
		}
		log.Printf("Login 500 Error: %v", err)
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, gin.H{
		"accessToken": token,
		"tokenType":   "Bearer",
	})
}

// GetMenus maneja GET /api/v1/auth/menus y devuelve los accesos del usuario logueado.
func (h *AuthHandler) GetMenus(c *gin.Context) {
	claims, ok := c.Get("auth_claims")
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "token no válido o no proporcionado")
		return
	}

	authClaims, ok := claims.(domain.AuthClaims)
	if !ok {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "error procesando claims")
		return
	}

	menus, err := h.service.GetMenus(c.Request.Context(), authClaims.IdEmpleado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, menus)
}

// GetPerfil maneja GET /api/v1/auth/perfil y devuelve los datos completos del operador (nombres, apellidos, foto).
func (h *AuthHandler) GetPerfil(c *gin.Context) {
	claims, ok := c.Get("auth_claims")
	if !ok {
		respondError(c, http.StatusUnauthorized, "UNAUTHORIZED", "token no válido o no proporcionado")
		return
	}

	authClaims, ok := claims.(domain.AuthClaims)
	if !ok {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "error procesando claims")
		return
	}

	perfil, err := h.service.GetUserProfile(c.Request.Context(), authClaims.IdEmpleado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, perfil)
}
