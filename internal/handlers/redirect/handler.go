package redirect

import (
	"fmt"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"

	"cutter-url-go/internal/models"
	"cutter-url-go/internal/repositories/link"
	"cutter-url-go/internal/vo"
)

const ValidationCode = http.StatusBadRequest

type Handler struct {
	r link.Repository
}

func NewHandler(r link.Repository) *Handler {
	return &Handler{r: r}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		h.getHandler(w, r)
	case http.MethodPost:
		h.postHandler(w, r)
	default:
		h.defaultHandler(w)
	}
}

func (h Handler) postHandler(w http.ResponseWriter, r *http.Request) {
	sl := h.getURI(w, r)
	if sl == nil {
		return
	}
	_, err := fmt.Fprint(w, sl.FullURL)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}

func (h Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	sl := h.getURI(w, r)
	if sl == nil {
		return
	}
	http.Redirect(w, r, string(sl.FullURL), http.StatusFound)
}

func (h Handler) defaultHandler(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h Handler) getURI(w http.ResponseWriter, r *http.Request) *models.ShortLink {
	// get body
	uri := vo.ShortURI(strings.Trim(r.URL.Path, "/"))

	// body validation
	if uri == "" {
		http.Error(w, "validation error", ValidationCode)
		return nil
	}

	sl, err := h.r.Get(uri)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.NotFound(w, r)
			return nil
		}
		http.Error(w, "find error", http.StatusInternalServerError)
		return nil
	}

	return sl
}
