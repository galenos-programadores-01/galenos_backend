package main

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"

	"github.com/galenos-pro/appointments-api/internal/config"
	"github.com/galenos-pro/appointments-api/internal/usecase"
)

func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	now := time.Now()
	claims := usecase.CustomClaims{
		IdEmpleado: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "admin",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(cfg.AuthSecret))
	if err != nil {
		fmt.Println("sign error:", err)
		return
	}
	fmt.Println(signed)
}
