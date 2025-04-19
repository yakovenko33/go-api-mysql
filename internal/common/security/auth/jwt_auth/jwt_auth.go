package jwt_auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	auth_errors "go-api-docker/internal/common/security/auth/errors"
	jwt_auth_repository "go-api-docker/internal/common/security/auth/infrastructure/jwt_auth_repository"
)

type JwtAuthManagerInterface interface {
	GenerateTokens(userData *UserData) (JwtTokens, error)
	VerifyToken(accessToken string) (string, error)
	RefreshTokens(refreshTokens string, userData *UserData) (JwtTokens, error)
	AddToBlackList(refreshTokens string) error
}

type JwtAuthManager struct {
	jwt_auth_repository jwt_auth_repository.JwtAuthRepositoryInterface
}

func NewJwtAuthManager(jwt_auth_repository jwt_auth_repository.JwtAuthRepositoryInterface) JwtAuthManagerInterface {
	return &JwtAuthManager{
		jwt_auth_repository: jwt_auth_repository,
	}
}

type JwtTokens struct {
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  int64
	RefreshTokenExpiry int64
}

type UserData struct {
	UserId    string
	UserAgent string
}

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func (m *JwtAuthManager) GenerateTokens(userData *UserData) (JwtTokens, error) {
	accessTokenExpiry := time.Now().UTC().Add(15 * time.Minute)
	accessTokenString, err := m.generateTokens(userData, &accessTokenExpiry)
	if err != nil {
		return JwtTokens{}, err
	}

	refreshTokenExpiry := time.Now().UTC().Add(30 * 24 * time.Hour)
	refreshTokenString, err := m.generateTokens(userData, &refreshTokenExpiry)
	if err != nil {
		return JwtTokens{}, err
	}

	jwtTokens := JwtTokens{
		AccessToken:        accessTokenString,
		RefreshToken:       refreshTokenString,
		AccessTokenExpiry:  accessTokenExpiry.Unix(),
		RefreshTokenExpiry: refreshTokenExpiry.Unix(),
	}
	_, err = m.addTokensInStore(userData, &jwtTokens, refreshTokenExpiry)

	return jwtTokens, err
}

func (m *JwtAuthManager) generateTokens(userData *UserData, expiresIn *time.Time) (string, error) {
	tokenClaims := jwt.MapClaims{
		"sub": userData.UserId,
		"exp": expiresIn.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (m *JwtAuthManager) addTokensInStore(userData *UserData, jwtTokens *JwtTokens, expiresIn time.Time) (string, error) {
	tokens := jwt_auth_repository.Tokens{
		ID:           uuid.New(),
		AccessToken:  jwtTokens.AccessToken,
		RefreshToken: jwtTokens.RefreshToken,
		UserAgent:    userData.UserAgent,
		UserId:       userData.UserId,
		BlackList:    0,
		ExpiresIn:    expiresIn,
		CreatedAt:    time.Now().UTC(),
	}
	result, err := m.jwt_auth_repository.AddTokens(&tokens)

	return result, err
}

func (m *JwtAuthManager) RefreshTokens(refreshTokenString string, userData *UserData) (JwtTokens, error) {
	token, err := m.jwt_auth_repository.GetTokensByRefreshToken(refreshTokenString)
	if token == nil {
		return JwtTokens{}, nil
	}
	if err != nil {
		return JwtTokens{}, err
	}

	userId, err := m.VerifyToken(refreshTokenString)
	if err != nil {
		return JwtTokens{}, auth_errors.NewAuthError(403, "refresh token invalid")
	}
	userData.UserId = userId

	newTokens, err := m.GenerateTokens(userData)
	if err != nil {
		return newTokens, err
	}

	err = m.jwt_auth_repository.AddToBlackList(token)
	if err != nil {
		return newTokens, err
	}

	return newTokens, nil
}

func (m *JwtAuthManager) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v", token.Header["alg"])
		}
		return string(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := claims["sub"].(string)
		return userID, nil
	}

	return "", auth_errors.NewAuthError(403, "token invalid")
}

func (m *JwtAuthManager) AddToBlackList(refreshTokens string) error {
	token, err := m.jwt_auth_repository.GetTokensByRefreshToken(refreshTokens)
	if token == nil {
		return nil
	}
	if err != nil {
		return err
	}

	err = m.jwt_auth_repository.AddToBlackList(token)
	if err != nil {
		return err
	}
	return nil
}
