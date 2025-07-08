package handlers

import (
	"context"
	"log"
	"medods_testcase/db/models"
	"medods_testcase/utils"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// @Summary		GetTokens
// @Tags			auth
// @Description	Получение пары токенов для указанного GUID
// @ID				get-tokens
// @Accept			json
// @Produce		json
// @Param			userId	query		string	true	"user guid" example(01ecfeac-5006-4af3-8961-fb5d3d7b142f)
// @Success		200		{object}	handlers.TokensResponse
// @Failure		default	{object}	handlers.ErrorResponse
// @Router			/api/auth/tokens [get]
func (h Handler) GetTokens(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("Запрос для GetTokens")
	w.Header().Set("Content-Type", "application/json")

	requestData := h.GetRequestData(r)

	// uGuid := r.URL.Query().Get("userId")
	// uAgent := r.UserAgent()
	// uIp := r.RemoteAddr

	queryGuid := uuid.FromStringOrNil(r.URL.Query().Get("userId"))

	log.Printf("Requestdata:\n%v\n", requestData)
	// log.Println(uGuid)
	// log.Println("UserAgent: ", uAgent)
	// log.Println("User IP: ", uIp)

	if queryGuid == uuid.Nil || requestData.userAgent == "" || requestData.userIp == "" {
		AssignResponse(map[string]string{"message": "Переданы некорректные данные"}, http.StatusBadRequest, w)
		return
	}

	tpg := utils.NewTokenPairGenerator(h.Config.AppAcsSecret, h.Config.AppRefSecret)

	res, err := tpg.GeneratePair(queryGuid.String(), requestData.userAgent, requestData.userIp)
	if res == nil || err != nil {
		log.Println("Ошибка генерации токена в GetTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка сервера"}, http.StatusInternalServerError, w)
		return
	} else {
		log.Printf("AccessToken:\n%v\nRefreshToken:\n%v\n", res.AccessToken, res.RefreshToken)
	}

	var token models.Token
	if err := h.DB.Where(models.Token{UserGuid: queryGuid}).FirstOrCreate(&token).Error; err != nil {
		log.Println("Ошибка получения токена в GetTokens", err)
		AssignResponse(map[string]string{"message": "Переданы некорректные данные"}, http.StatusBadRequest, w)
		return
	}
	token.Refresh_token = tpg.HashedRefreshToken
	token.Last_fingerprint = requestData.userAgent
	token.Last_ip = requestData.userIp
	token.PairId = tpg.PairId

	if err := h.DB.Save(&token).Error; err != nil {
		log.Println("Ошибка сохранения записи в GetTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка сервера"}, http.StatusInternalServerError, w)
		return
	}

	var user models.User
	if err := h.DB.First(&user, "guid = ?", queryGuid).Error; err != nil {
		log.Println("Ошибка чтения записей в GetTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка сервера"}, http.StatusInternalServerError, w)
		return
	}

	refreshCookie, accessCookie := SetTokenCookies(tpg)

	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)

	AssignResponse(map[string]string{"message": "Токены получены", "refresh_token": tpg.EncodedRefreshToken, "access_token": tpg.AccessToken}, http.StatusOK, w)
}

func GetRequestTokens(ctx context.Context) (string, string) {
	refreshToken, _ := ctx.Value(RefKey).(string)
	accessToken, _ := ctx.Value(AcsKey).(string)
	// accessToken := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	return refreshToken, accessToken
}

func GetTokenCookies(r *http.Request) (*http.Cookie, *http.Cookie) {
	accessCookie, err := r.Cookie("access_token")
	if err != nil {
		log.Println("Ошибка получения accessCookie", err)
		return nil, nil
	}

	refreshCookie, err := r.Cookie("access_token")
	if err != nil {
		log.Println("Ошибка получения refreshCookie", err)
		return nil, nil
	}

	return accessCookie, refreshCookie
}

func SetTokenCookies(pg *utils.TokenPairGenerator) (*http.Cookie, *http.Cookie) {
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Domain:   "localhost",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		MaxAge:   int(60 * 60 * 24 * 30),
		Value:    pg.EncodedRefreshToken,
	}

	accessCookie := &http.Cookie{
		Name:     "access_token",
		Domain:   "localhost",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 30),
		MaxAge:   int(60 * 30),
		Value:    pg.AccessToken,
	}

	return refreshCookie, accessCookie
}

func ResetTokenCookies() (*http.Cookie, *http.Cookie) {
	refreshCookie := &http.Cookie{
		Name:   "refresh_token",
		Domain: "localhost",
		Path:   "/",
		MaxAge: -1,
		Value:  "",
	}

	accessCookie := &http.Cookie{
		Name:   "access_token",
		Domain: "localhost",
		Path:   "/",
		MaxAge: -1,
		Value:  "",
	}
	return refreshCookie, accessCookie
}
