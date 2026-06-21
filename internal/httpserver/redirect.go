package httpserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mereska0/cliplink/internal/domain"
	"github.com/mereska0/cliplink/internal/service"
)

type RedirectHandler struct {
	linkService *service.LinkService
}

func NewRedirectHandler(linkServ *service.LinkService) *RedirectHandler {
	return &RedirectHandler{linkService: linkServ}
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}
	link, err := h.linkService.RegisterClick(r.Context(), code)

	if err != nil {
		if errors.Is(err, domain.ErrLinkNotFound) {
			http.Error(w, "link not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, domain.ErrDeletedLink) {
			http.Error(w, "link is deleted", http.StatusGone)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, link.OriginalURL, http.StatusFound)
}
