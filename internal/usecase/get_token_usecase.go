package usecase

import (
	"fmt"
	"net/http"
	"strings"
)

func (uc *UseCase) GetBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// Проверяем, что заголовок начинается с "Bearer "
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("authorization header does not start with 'Bearer '")
	}

	// Извлекаем токен, удаляя префикс "Bearer "
	token := strings.TrimPrefix(authHeader, bearerPrefix)

	return token, nil
}
