package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type authUseCase struct {
	repo   output.AuthRepository
	secret []byte
	ttl    time.Duration
}

type CustomClaims struct {
	IdEmpleado int `json:"idEmpleado"`
	jwt.RegisteredClaims
}

// NewAuthUseCase construye el caso de uso de autenticación, conectado a la BD.
func NewAuthUseCase(repo output.AuthRepository, secret string, ttl time.Duration) input.AuthService {
	return &authUseCase{
		repo:   repo,
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (string, error) {
	// Llamar a la BD para validar
	idEmpleado, err := uc.repo.Login(ctx, username, password)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := CustomClaims{
		IdEmpleado: idEmpleado,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(uc.secret)
	if err != nil {
		return "", fmt.Errorf("signing jwt: %w", err)
	}

	return signed, nil
}

func (uc *authUseCase) ValidateToken(tokenString string) (domain.AuthClaims, error) {
	var claims CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return uc.secret, nil
	})
	if err != nil || !token.Valid {
		return domain.AuthClaims{}, domain.ErrInvalidToken
	}

	return domain.AuthClaims{
		Subject:    claims.Subject,
		IdEmpleado: claims.IdEmpleado,
	}, nil
}

func (uc *authUseCase) GetMenus(ctx context.Context, idEmpleado int) (domain.AuthMenus, error) {
	menus, err := uc.repo.GetMenus(ctx, idEmpleado)
	if err != nil {
		return domain.AuthMenus{}, err
	}

	permisos, err := uc.repo.GetMenuPermisos(ctx, idEmpleado)
	if err != nil {
		return domain.AuthMenus{}, err
	}

	return domain.AuthMenus{
		Menus:    menus,
		Permisos: permisos,
	}, nil
}

func (uc *authUseCase) GetUserProfile(ctx context.Context, idEmpleado int) (domain.UserProfile, error) {
	return uc.repo.GetUserProfile(ctx, idEmpleado)
}
