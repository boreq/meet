package meet

import (
	"net/http"

	"github.com/boreq/meet/application"
	"github.com/boreq/meet/application/meet"
	"github.com/boreq/meet/domain"
	"github.com/boreq/meet/internal/logging"
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
	h.router.Get("/meetings/{meetingName}/websocket", h.meetingWebsocket)
}

func (h *Handler) meetingWebsocket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	h.log.Debug("connecting to the meeting websocket")

	meetingName, err := domain.NewMeetingName(chi.URLParam(r, "meetingName"))
	if err != nil {
		h.log.Debug("invalid meeting name", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// todo why is this needed on localhost?
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Debug("upgrade failed", "err", err)
		return
	}

	cmd := meet.JoinMeeting{
		MeetingName: meetingName,
		Client:      NewClient(conn),
	}

	if err := h.app.Meet.JoinMeeting.Execute(ctx, cmd); err != nil {
		h.log.Debug("join meeting error", err, "err")
	}
}
