package handlers

import (
	"log"
	"medods_testcase/utils"
	"net/http"

	"github.com/gofrs/uuid"
)

// @Summary		GetGuid
// @Security		ApiTokenAuth
// @Tags			auth
// @Description	Получение GUID пользователя
// @ID				get-guid
// @Accept			json
// @Produce		json
// @Param			input	body		handlers.TokensRequest	true	"user tokens"
// @Success		200		{object}	handlers.GuidResponse
// @Failure		default	{object}	handlers.ErrorResponse
// @Router			/api/auth/user [get]
func (h Handler) GetGuid(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("Запрос для GetGuid")
	w.Header().Set("Content-Type", "application/json")

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
		AssignResponse(map[string]string{"message": "Ошибка получения GUID"}, http.StatusBadRequest, w)
		return
	} else {
		queryGuid = uuid.FromStringOrNil(userGuid)
	}

	AssignResponse(map[string]string{"message": "GUID получен", "GUID": queryGuid.String()}, http.StatusOK, w)
}
