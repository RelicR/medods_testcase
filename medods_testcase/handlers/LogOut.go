package handlers

import (
	"log"
	"net/http"
)

// LogOut godoc
//
// @Summary		LogOut
// @Tags			auth
// @Description	Деавторизация пользователя
// @ID				log-out
// @Accept			json
// @Produce		json
// @Success		200		{object}	handlers.MessageResponse
// @Router			/api/auth/logout [post]
func (h Handler) LogOut(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("Запрос для LogOut")
	w.Header().Set("Content-Type", "application/json")

	refreshCookie, accessCookie := ResetTokenCookies()

	http.SetCookie(w, refreshCookie)
	http.SetCookie(w, accessCookie)
	AssignResponse(map[string]string{"message": "Успешный выход из аккаунта"}, http.StatusOK, w)
}
