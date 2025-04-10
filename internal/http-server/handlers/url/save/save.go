package save

import (
	"errors"
	resp "github.com/UncleWeb1992/Url-SHortner/internal/lib/api/response"
	"github.com/UncleWeb1992/Url-SHortner/internal/lib/logger/sl"
	"github.com/UncleWeb1992/Url-SHortner/internal/lib/utils/random"
	"github.com/UncleWeb1992/Url-SHortner/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias, omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias, omitempty"`
}

type UrlSaver interface {
	SaveUrl(url string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver, aliasLength int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const repo = "handlers.url.save.New"

		log = log.With(
			slog.String("repo", repo),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("cannot decode request", sl.Err(err))
			render.JSON(w, r, resp.Error("cannot decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("cannot validate request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias

		if alias == "" {
			alias = random.GetRandomString(aliasLength)
		}

		id, err := urlSaver.SaveUrl(req.URL, alias)

		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("url already exists"))
			return
		}

		if err != nil {
			log.Error("cannot save url", sl.Err(err))

			render.JSON(w, r, resp.Error("cannot save url"))
			return
		}

		log.Info("url saved", slog.Int64("id", id), slog.String("alias", alias))

		responseOk(w, r, alias)
	}
}

func responseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.Ok(),
		Alias:    alias,
	})
}
