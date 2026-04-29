package websocket

import (
	"log"
	"net/http"
)

func (s *Service) HandleWS(
	w http.ResponseWriter,
	r *http.Request,
) {

	roomID := r.URL.Query().Get("room_id")
	userID := r.URL.Query().Get("user_id")

	if roomID == "" || userID == "" {
		http.Error(w, "missing params", 400)
		return
	}

	err := s.ConnectUser(
		w,
		r,
		roomID,
		userID,
	)

	if err != nil {
		http.Error(w, err.Error(), 403)
	}
	log.Println("WS connected:", userID, roomID)
}
