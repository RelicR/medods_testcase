package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// @Summary		Webhook
// @Tags			utils
// @Description	Пример обработки вебхука
// @ID				webhook
// @Accept			json
// @Produce		json
// @Param			input	body		handlers.Payload	true	"user guid" example(01ecfeac-5006-4af3-8961-fb5d3d7b142f)
// @Success		200		{object}	handlers.WebhookResponse "status"
// @Failure		default	{object}	handlers.WebhookResponse "status"
// @Router			/webhook [post]
func (h Handler) Webhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	log.Println("Событие webhook")
	w.Header().Set("Content-Type", "application/json")

	payload := &Payload{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Ошибка чтения тела запроса в Webhook", err)
		AssignResponse(map[string]string{"status": "failure"}, http.StatusInternalServerError, w)
		return
	}

	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("Ошибка преобразования тела запроса в Webhook", err)
		AssignResponse(map[string]string{"status": "failure"}, http.StatusInternalServerError, w)
		return
	}

	log.Printf("Событие webhook:\n%v", payload)

	AssignResponse(map[string]string{"status": "ok"}, http.StatusOK, w)
}

type Payload struct {
	UserGuid  string `json:"guid" example:"01ecfeac-5006-4af3-8961-fb5d3d7b142f"`
	EventType uint16 `json:"event" example:"111"`
	EventInfo string `json:"info" example:"192.168.0.1"`
}

type WebhookResponse struct {
	Status string `json:"status" example:"ok"`
}
