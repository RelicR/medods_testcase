package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AccessTokenClaims struct {
	Fp     string `json:"fingerprint"`
	Ip     string `json:"ip"`
	PairId string
	jwt.RegisteredClaims
}

type TokenPairGenerator struct {
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
	AccessTokenSecret   []byte
	RefreshTokenSecret  []byte
	HashedRefreshToken  string
	EncodedRefreshToken string
	PairId              string
}

func NewTokenPairGenerator(accessSecret, refreshSecret string) *TokenPairGenerator {
	return &TokenPairGenerator{
		AccessTokenSecret:  []byte(accessSecret),
		RefreshTokenSecret: []byte(refreshSecret),
	}
}

func GeneratePairId(length int) (string, error) {
	fillingMas := make([]byte, length)

	if _, err := rand.Read(fillingMas); err != nil {
		log.Panic("Ошибка генерации ключа пары токенов в GeneratePairId")
		return "", err
	}

	masHasher := sha256.New()
	masHasher.Write(fillingMas)

	return base64.URLEncoding.EncodeToString(masHasher.Sum(nil))[:length], nil
}

func (pg *TokenPairGenerator) GeneratePair(userId, userAgent, userIp string) (*TokenPairGenerator, error) {
	if pairId, err := GeneratePairId(16); err != nil {
		return nil, err
	} else {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessTokenClaims{
			userAgent,
			userIp,
			pairId,
			jwt.RegisteredClaims{
				Subject:   userId,
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
				Issuer:    "auth.localhost:8080",
				Audience:  []string{"localhost:8080"},
			},
		})

		accessSigned, err := accessToken.SignedString(pg.AccessTokenSecret)
		if err != nil {
			return nil, err
		}

		refreshToken, err := uuid.NewV7()
		if err != nil {
			return nil, err
		}

		pg.AccessToken = accessSigned
		pg.RefreshToken = refreshToken.String()
		pg.PairId = pairId
		pg.EncodedRefreshToken = base64.StdEncoding.Strict().EncodeToString([]byte(refreshToken.String()))

		if res, err := bcrypt.GenerateFromPassword([]byte(pg.RefreshToken), bcrypt.DefaultCost); err != nil {
			return nil, err
		} else {
			pg.HashedRefreshToken = string(res)
		}

		return pg, nil
	}

}

func (pg *TokenPairGenerator) VerifyPair() bool {
	accessClaims, err := pg.ParseAccessToken()
	if err != nil {
		log.Println("Ошибка токена доступа", err)
		return false
	}
	if accessClaims.PairId != pg.PairId {
		log.Printf("Недействительная пара токенов:\nRefresh %s\nAccess %s\n", pg.PairId, accessClaims.PairId)
		return false
	}
	return true
}

func (pg *TokenPairGenerator) ParseAccessToken() (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(pg.AccessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неверный алгоритм подписи %v", token.Method.Alg())
		}
		return pg.AccessTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	log.Printf("Claims:\n%v", token.Claims.(*AccessTokenClaims))
	log.Printf("IsValid:\n%v", token.Valid)
	return nil, fmt.Errorf("недействительный токен")
}

func (pg *TokenPairGenerator) GetUserGuid() (string, error) {
	if claims, err := pg.ParseAccessToken(); err != nil {
		return "", err
	} else {
		log.Printf("Claims is:\n%v", claims)
		return claims.GetSubject()
	}
}
