package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"medods_testcase/db"
	"medods_testcase/handlers"
	"medods_testcase/utils"
	"net/http"

	_ "medods_testcase/docs"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Тестовое задание для API сервиса аутентификации
//	@version		1.0
//	@description	API server for auth tokens
//	@host			localhost:8081
//	@BasePath		/

// @securitydefinitions.apikey	ApiTokenAuth
// @in							header
// @name						Authorization
func main() {
	defer log.Panicln("Сервис завершил работу")

	config := utils.NewConfig()

	db := db.NewDbCon()
	hd := handlers.NewHandler(db, config)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/auth/tokens", hd.GetTokens)
	mux.HandleFunc("POST /api/auth/tokens/refresh", ProtectRoute(hd.RefreshTokens))
	mux.HandleFunc("GET /api/auth/user", ProtectRoute(hd.GetGuid))
	mux.HandleFunc("POST /api/auth/logout", hd.LogOut)

	mux.HandleFunc("POST /webhook", hd.Webhook)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.AppHost, config.AppPort), mux)
	log.Printf("Сервис запущен на %s:%d", config.AppHost, config.AppPort)
}

func ProtectRoute(hf func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Header.Get("Authorization")
		log.Printf("Auth header is:%v", userId)
		var payload struct {
			RefreshToken string `json:"refresh_token"`
			AccessToken  string `json:"access_token"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Ошибка чтения тела запроса в ProtectRoute", err)
			handlers.AssignResponse(map[string]string{"status": "failure"}, http.StatusUnauthorized, w)
			return
		}

		log.Printf("Body is:\n%v", body)

		if err := json.Unmarshal(body, &payload); err != nil {
			log.Println("Ошибка преобразования тела запроса в ProtectRoute", err)
			handlers.AssignResponse(map[string]string{"status": "failure"}, http.StatusUnauthorized, w)
			return
		}

		decodedRefreshToken, err := base64.StdEncoding.Strict().DecodeString(payload.RefreshToken)
		log.Printf("Decoded b64 refresh:\n%v\n", string(decodedRefreshToken))
		if err != nil {
			log.Println("Ошибка преобразования refresh в ProtectRoute", err)
			handlers.AssignResponse(map[string]string{"status": "failure"}, http.StatusUnauthorized, w)
			return
		} else {
			r = r.WithContext(context.WithValue(r.Context(), handlers.RefKey, string(decodedRefreshToken)))
			r = r.WithContext(context.WithValue(r.Context(), handlers.AcsKey, payload.AccessToken))
			log.Printf("Payload is:%v", payload)
			hf(w, r)
		}
	}
}

func init() {
	if err := godotenv.Load("app.env"); err != nil {
		log.Print("Не найден файл .env")
	}
}
