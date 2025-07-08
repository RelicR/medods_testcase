package handlers

import (
	"encoding/json"
	"log"
	"medods_testcase/utils"
	"net/http"

	"gorm.io/gorm"
)

type Handler struct {
	DB     *gorm.DB
	Config *utils.Config
}

func NewHandler(db *gorm.DB, conf *utils.Config) Handler {
	log.Println("Создан новый обработчик")
	return Handler{db, conf}
}

func AssignResponse(data map[string]string, status int, writer http.ResponseWriter) {
	writer.WriteHeader(status)

	json.NewEncoder(writer).Encode(data)
}

type RequestData struct {
	userAgent string
	userIp    string
}

func (h Handler) GetRequestData(r *http.Request) *RequestData {
	rd := &RequestData{}

	rd.userAgent = r.UserAgent()
	rd.userIp = r.RemoteAddr

	return rd
}

type ctxKey int

const AcsKey ctxKey = 0
const RefKey ctxKey = 1

type TokensRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJmaW5nZXJwcmludCI6IlBvc3RtYW5SdW50aW1lLzcuNDQuMSIsImlwIjoiMTcyLjE4LjAuMTo0NDU2OCIsIlBhaXJJZCI6IklFbjdHZEZ2dm1BNXo3M3giLCJpc3MiOiJhdXRoLmxvY2FsaG9zdDo4MDgwIiwic3ViIjoiMDFlY2ZlYWMtNTAwNi00YWYzLTg5NjEtZmI1ZDNkN2IxNDJmIiwiYXVkIjpbImxvY2FsaG9zdDo4MDgwIl0sImV4cCI6MTc1MTk5Mzg5OCwiaWF0IjoxNzUxOTkyMDk4fQ.veQb44AcPIr1J0oxvhyaaCbrzo6gV-sMLM1xtRJDN_p3tOfR-lSwJTR8K5oqL-CbfJ9j_OTFkpHlEUeN3eFmVA"`
	AccessToken  string `json:"access_token" example:"MDE5N2VhZGQtZGQxMC03MDE4LTg4ZWYtODc5N2EzNGIxMjk3"`
}

type TokensResponse struct {
	Message      string `json:"message" example:"Токены обновлены"`
	RefreshToken string `json:"refresh_token" example:"MDE5N2VhZGUtMDQ3NC03MDE4LWJlMmYtMGUxMzliMTY0M2Ex"`
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJmaW5nZXJwcmludCI6IlBvc3RtYW5SdW50aW1lLzcuNDQuMSIsImlwIjoiMTcyLjE4LjAuMTo0NDU2OCIsIlBhaXJJZCI6IkZyRW5fX0o4UTlSc1djcVUiLCJpc3MiOiJhdXRoLmxvY2FsaG9zdDo4MDgwIiwic3ViIjoiMDFlY2ZlYWMtNTAwNi00YWYzLTg5NjEtZmI1ZDNkN2IxNDJmIiwiYXVkIjpbImxvY2FsaG9zdDo4MDgwIl0sImV4cCI6MTc1MTk5MzkwOCwiaWF0IjoxNzUxOTkyMTA4fQ.YA-Cp0iNzeWwX4Z7zRllpq8tIOpE5KIJWuFpU5u49sRtPdHkLFGLAul744APJe5UZ0fnUuc8yfWpCLiNRCtogw"`
}

type GuidResponse struct {
	Message string `json:"message" example:"GUID получен"`
	GUID    string `json:"guid" example:"01ecfeac-5006-4af3-8961-fb5d3d7b142f"`
}

type MessageResponse struct {
	Message string `json:"message" example:"Успешный выход из аккаунта"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"Переданы некорректные данные"`
}
