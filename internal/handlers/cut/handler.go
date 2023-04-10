package cut

import (
	"errors"
	"net/http"
	"net/url"

	"go.mongodb.org/mongo-driver/mongo"

	"cutter-url-go/internal/models"
	"cutter-url-go/internal/repositories/link"
	"cutter-url-go/internal/vo"
)

const ValidationCode = http.StatusBadRequest
const AlreadyUsedCode = http.StatusConflict

type Handler struct {
	r link.Repository
}

func NewHandler(r link.Repository) *Handler {
	return &Handler{r: r}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		h.postHandler(w, r)
	default:
		h.defaultHandler(w)
	}
}

func (h Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	// get body
	shortURI := r.PostFormValue("short_uri")
	fullURI := r.PostFormValue("full_uri")

	// body validation
	if shortURI == "" || fullURI == "" {
		http.Error(w, "validation error", ValidationCode)
		return
	}
	if _, err := url.ParseRequestURI(fullURI); err != nil {
		http.Error(w, "validation error", ValidationCode)
		return
	}

	sl := models.ShortLink{
		FullURL:   vo.FullURI(fullURI),
		ShortLink: vo.ShortURI(shortURI),
	}
	err := h.r.Insert(&sl)
	if err != nil {
		e := mongo.WriteException{}
		if errors.As(err, &e) {
			for _, we := range e.WriteErrors {
				if we.Code == 11000 {
					http.Error(w, "short URI already used", AlreadyUsedCode)
					return
				}
			}
		}
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handler) defaultHandler(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
