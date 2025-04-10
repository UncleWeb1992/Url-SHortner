package redirect

import (
	resp "github.com/UncleWeb1992/Url-SHortner/internal/lib/api/response"
	"github.com/UncleWeb1992/Url-SHortner/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	Url string `json:"url"`
}

type Redirect interface {
	RedirectByAlias(alias string) (string, error)
}

func New(log *slog.Logger, RedirectByAlias Redirect) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const repo = "handlers.redirect.New"

		log = log.With(
			slog.String("repo", repo),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			render.JSON(w, r, resp.Error("alias is empty"))
			return
		}

		url, err := RedirectByAlias.RedirectByAlias(alias)

		if err != nil {
			log.Error("cannot redirect", sl.Err(err))
			render.JSON(w, r, resp.Error("cannot redirect"))
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}
