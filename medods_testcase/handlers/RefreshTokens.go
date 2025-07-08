package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"medods_testcase/db/models"
	"medods_testcase/utils"
	"net/http"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// @Summary		RefreshTokens
// @Security		ApiTokenAuth
// @Tags			auth
// @Description	Обновление токенов
// @ID				refresh-tokens
// @Accept			json
// @Produce		json
// @Param			input	body		handlers.TokensRequest	true	"user tokens"
// @Success		200		{object}	handlers.TokensResponse
// @Failure		default	{object}	handlers.ErrorResponse
// @Router			/api/auth/tokens/refresh [post]
func (h Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("Запрос для RefreshTokens")
	w.Header().Set("Content-Type", "application/json")

	requestData := h.GetRequestData(r)
	var queryGuid uuid.UUID

	tpg := utils.NewTokenPairGenerator(h.Config.AppAcsSecret, h.Config.AppRefSecret)

	// if accessCookie, refreshCookie := GetTokenCookies(r); accessCookie == nil || refreshCookie == nil {
	// 	AssignResponse(map[string]string{"message": "Переданы некорректные данные"}, http.StatusBadRequest, w)
	// 	return
	// } else {
	// 	tpg.AccessToken = accessCookie.Value
	// 	tpg.RefreshToken = refreshCookie.Value
	// }

	if refreshToken, accessToken := GetRequestTokens(r.Context()); refreshToken == "" || accessToken == "" {
		AssignResponse(map[string]string{"message": "Переданы некорректные данные"}, http.StatusBadRequest, w)
		return
	} else {
		tpg.AccessToken = accessToken
		tpg.RefreshToken = refreshToken
	}

	if userGuid, err := tpg.GetUserGuid(); err != nil {
		log.Println("Ошибка получения GUID в RefreshTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка обновления токенов доступа"}, http.StatusBadRequest, w)
		return
	} else {
		queryGuid = uuid.FromStringOrNil(userGuid)
	}

	var refreshToken models.Token
	if err := h.DB.Where(models.Token{UserGuid: queryGuid}).First(&refreshToken).Error; err != nil {
		log.Println("Ошибка получения токена в RefreshTokens", err)
		AssignResponse(map[string]string{"message": "Переданы некорректные данные"}, http.StatusBadRequest, w)
		return
	}

	compareRefresh := bcrypt.CompareHashAndPassword([]byte(refreshToken.Refresh_token), []byte(tpg.RefreshToken))
	tpg.PairId = refreshToken.PairId
	if valid := tpg.VerifyPair(); !valid || refreshToken.Last_fingerprint != requestData.userAgent || compareRefresh != nil {
		log.Println("Ошибка подтверждения refresh токена в RefreshTokens")

		refreshCookie, accessCookie := ResetTokenCookies()
		http.SetCookie(w, refreshCookie)
		http.SetCookie(w, accessCookie)
		AssignResponse(map[string]string{"message": "Ошибка обновления токенов доступа"}, http.StatusBadRequest, w)
		return
	}

	if requestData.userIp != refreshToken.Last_ip {
		SendWebhook(queryGuid.String(), requestData.userIp, 111)
	}

	res, err := tpg.GeneratePair(queryGuid.String(), requestData.userAgent, requestData.userIp)
	if res == nil || err != nil {
		log.Println("Ошибка генерации токена в RefreshTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка сервера"}, http.StatusInternalServerError, w)
		return
	} else {
		log.Printf("AccessToken:\n%v\nRefreshToken:\n%v\n", res.AccessToken, res.RefreshToken)
	}

	refreshToken.Refresh_token = tpg.HashedRefreshToken
	refreshToken.Last_fingerprint = requestData.userAgent
	refreshToken.Last_ip = requestData.userIp
	refreshToken.PairId = tpg.PairId

	if err := h.DB.Save(&refreshToken).Error; err != nil {
		log.Println("Ошибка чтения записей в GetTokens", err)
		AssignResponse(map[string]string{"message": "Ошибка сервера"}, http.StatusInternalServerError, w)
		return
	}

	refreshCookie, accessCookie := SetTokenCookies(tpg)

	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)

	AssignResponse(map[string]string{"message": "Токены обновлены", "refresh_token": tpg.EncodedRefreshToken, "access_token": tpg.AccessToken}, http.StatusOK, w)
}

func SendWebhook(userId, newLogonIp string, eventType uint16) {
	webhookEndpoint := "http://localhost:8080/webhook"

	payload := map[string]interface{}{
		"guid":  userId,
		"event": eventType,
		"info":  newLogonIp,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Println("Ошибка сериализации", err)
		return
	}

	req, err := http.NewRequest("POST", webhookEndpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Println("Ошибка формирования запроса", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Ошибка отправки запроса", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Ошибка обработки вебхука", err)
		return
	}

	log.Println("Вебхук обработан", err)
}
