package helper

import (
	"errors"
	"fmt"
	"simple-crud-clean-architecture/internal/entity"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type TokenHelper struct {
	secret     *Secret
	expireTime *ExpireTime
}

type Secret struct {
	accessSecretByte  []byte
	refreshSecretByte []byte
	keySecret         []byte
}

type ExpireTime struct {
	AccessToken         time.Duration
	RefreshToken        time.Duration
	EmployeeAccessToken time.Duration
}

func NewTokenHelper(config *viper.Viper, log *logrus.Logger) *TokenHelper {
	accessSecret := config.GetString("app.access_secret")
	refreshSecret := config.GetString("app.refresh_secret")
	keySecret := config.GetString("app.key_secret")
	secret := &Secret{
		accessSecretByte:  []byte(accessSecret),
		refreshSecretByte: []byte(refreshSecret),
		keySecret:         []byte(keySecret),
	}

	expireTime := &ExpireTime{
		AccessToken:         time.Duration(config.GetInt("app.access_token_expire_in_minute")),
		RefreshToken:        time.Duration(config.GetInt("app.refresh_token_expire_in_day")),
		EmployeeAccessToken: time.Duration(config.GetInt("app.employee_access_token_expire_in_day")),
	}

	return &TokenHelper{
		secret:     secret,
		expireTime: expireTime,
	}
}

func (c *TokenHelper) GenerateEmployeeAccessToken(employeeUUID uuid.UUID) (*entity.EmployeeAccessToken, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * c.expireTime.EmployeeAccessToken)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["employee_uuid"] = employeeUUID
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	stringToken, err := token.SignedString(c.secret.accessSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.EmployeeAccessToken{
		EmployeeUUID: employeeUUID,
		Token:        stringToken,
		ExpiresAt:    expirationTime,
	}, nil
}

func (c *TokenHelper) GenerateAccessToken(userUUID uuid.UUID) (*entity.AccessToken, error) {
	expirationTime := time.Now().Add(time.Minute * c.expireTime.AccessToken)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_uuid"] = userUUID
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	stringToken, err := token.SignedString(c.secret.accessSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.AccessToken{
		UserUUID:  userUUID,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *TokenHelper) GenerateRefreshToken(userUUID uuid.UUID) (*entity.RefreshToken, error) {
	expirationTime := time.Now().Add(time.Hour * 24 * c.expireTime.RefreshToken)

	claims := jwt.MapClaims{}
	claims["user_uuid"] = userUUID
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString(c.secret.refreshSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.RefreshToken{
		UserUUID:  userUUID,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *TokenHelper) VerifyEmployeeAccessToken(token string) (*entity.EmployeeAccessToken, error) {

	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.secret.accessSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenDetails := &entity.EmployeeAccessToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		employeeIDStr, ok := claims["employee_uuid"].(string)
		if !ok {
			fmt.Println("employee_uuid not a string")
			return nil, errors.New("invalid token claims")
		}

		employeeUUID, err := uuid.Parse(employeeIDStr)
		if err != nil {
			fmt.Println("failed to parse uuid:", err)
			return nil, errors.New("invalid token claims")
		}

		accessTokenDetails.EmployeeUUID = employeeUUID
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		accessTokenDetails.ExpiresAt = expiresAt
	}

	return accessTokenDetails, nil

}

func (c *TokenHelper) VerifyAccessToken(token string) (*entity.AccessToken, error) {

	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.secret.accessSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenDetails := &entity.AccessToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		userIDStr, ok := claims["user_uuid"].(string)
		if !ok {
			fmt.Println("user_uuid not a string")
			return nil, errors.New("invalid token claims")
		}

		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			fmt.Println("failed to parse uuid:", err)
			return nil, errors.New("invalid token claims")
		}

		accessTokenDetails.UserUUID = userUUID
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			return nil, errors.New("invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		accessTokenDetails.ExpiresAt = expiresAt
	}

	return accessTokenDetails, nil

}

func (c *TokenHelper) VerifyRefreshToken(token string) (*entity.RefreshToken, error) {

	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.secret.refreshSecretByte, nil
	})
	if err != nil {
		fmt.Println("Error in parsing token:", err)
		return nil, err
	}

	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if !ok || !tokenClaims.Valid {
		return nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_uuid"].(string)
	if !ok {
		fmt.Println("user_uuid not a string")
		return nil, errors.New("invalid token claims")
	}

	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		fmt.Println("failed to parse uuid:", err)
		return nil, errors.New("invalid token claims")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid exp in token claims")
	}
	expiresAt := time.Unix(int64(expFloat), 0)

	refreshTokenDetails := &entity.RefreshToken{
		UserUUID:  userUUID,
		ExpiresAt: expiresAt,
	}

	return refreshTokenDetails, nil
}

func (c *TokenHelper) Encrypt(plaintext string) (string, error) {
	key := c.secret.keySecret
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

func (c *TokenHelper) Decrypt(ciphertext string) (string, error) {
	key := c.secret.keySecret
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fiber.ErrBadRequest
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fiber.ErrBadRequest
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", fiber.ErrBadRequest
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}
