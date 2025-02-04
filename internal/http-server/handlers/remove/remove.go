package remove

import (
	"errors"
	"net/http"
	"pet-project/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"

	resp "pet-project/internal/lib/api/response"

	"pet-project/internal/storage"
)

//go:generate go run github.com/vektra/mockery/v2@v2.51.1 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) (int64, error)
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.remove.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		rowsDeleted, err := urlRemover.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", slog.String("alias", alias))
			render.JSON(w, r, resp.Error("not found"))
			return
		}
		if err != nil {
			log.Error("failed to delete url", sl.Err(err), slog.String("alias", alias))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url deleted successfully", slog.Int64("rows_deleted", rowsDeleted))
		render.JSON(w, r, resp.OK())
	}
}
