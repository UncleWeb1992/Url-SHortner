package delete

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
	resp.Response
	Id int64
}

type UrlDeleted interface {
	DeleteUrl(alias string) (int64, error)
}

func New(log *slog.Logger, urlDeleted UrlDeleted) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const repo = "handlers.url.delete.New"

		log = log.With(
			slog.String("repo", repo),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		id, err := urlDeleted.DeleteUrl(alias)

		if err != nil {
			log.Error("cannot delete url", sl.Err(err))
			render.JSON(w, r, resp.Error("cannot delete url"))
			return
		}

		render.JSON(w, r, Response{Id: id, Response: resp.Ok()})
	}
}
