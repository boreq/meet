package meet

import (
	"encoding/json"
	"net/http"

	"github.com/boreq/meet/adapters/meet"
	"github.com/boreq/meet/application"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
	"github.com/boreq/rest"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type Handler struct {
	app    *application.Application
	router *chi.Mux
	log    logging.Logger
}

func NewHandler(app *application.Application) *Handler {
	h := &Handler{
		app:    app,
		router: chi.NewRouter(),
		log:    logging.New("ports/http/meet"),
	}

	h.registerHandlers()

	return h
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(rw, req)
}

func (h *Handler) registerHandlers() {
	h.router.Post("/meetings/{meetingName}/websocket", h.meeting)
	h.router.Post("/meetings/{meetingName}/sdp", rest.Wrap(h.sdp))
}

func (h *Handler) meeting(w http.ResponseWriter, r *http.Request) {
	_, err := domain.NewMeetingName(chi.URLParam(r, "meetingName"))
	if err != nil {
		h.log.Debug("invalid meeting name", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	_, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Debug("upgrade failed", "err", err)
		return
	}

}

var meeting = meet.NewWebRTCMeting()

func (h *Handler) sdp(r *http.Request) rest.RestResponse {
	meetingName, err := domain.NewMeetingName(chi.URLParam(r, "meetingName"))
	if err != nil {
		h.log.Debug("invalid meeting name", "err", err)
		return rest.ErrBadRequest
	}

	var request JoinMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.log.Debug("json decoding failed", "err", err)
		return rest.ErrBadRequest
	}

	h.log.Debug("sdp received", "meetingName", meetingName)

	member, err := meeting.Join(request.Sdp)
	if err != nil {
		h.log.Debug("sdp failure", "err", err)
		return rest.ErrInternalServerError
	}

	answer, err := member.Answer()
	if err != nil {
		h.log.Debug("answer failure", "err", err)
		return rest.ErrInternalServerError
	}

	return rest.NewResponse(JoinMeetingResponse{answer})
}
