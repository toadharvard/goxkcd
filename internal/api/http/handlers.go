package http

import (
	"context"
	"net/http"
	"strconv"

	"log/slog"

	"github.com/toadharvard/goxkcd/internal/usecase/buildIndex"
	"github.com/toadharvard/goxkcd/internal/usecase/countComics"
	"github.com/toadharvard/goxkcd/internal/usecase/downloadComics"
	"github.com/toadharvard/goxkcd/internal/usecase/suggestComix"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

func PingHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = Encode(w, r, http.StatusOK, "pong")
	}
}

type SuggestComixPicsResponse struct {
	Pics []string `json:"pics"`
}

type SuggestComixPicsRequest struct {
	language iso6391.ISOCode6391
	query    string
	limit    int
}

func NewSuggestComixPicsRequest(r *http.Request) (req SuggestComixPicsRequest, err error) {
	language, err := iso6391.NewLanguage(
		r.URL.Query().Get("language"),
	)

	if err != nil {
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		return
	}

	return SuggestComixPicsRequest{
		language: language,
		query:    query,
		limit:    limit,
	}, nil
}

func SuggestComixPicsHandler(ctx context.Context, suggestComix *suggestComix.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := NewSuggestComixPicsRequest(r)
		if err != nil {
			err = Encode(w, r, http.StatusBadRequest, err.Error())
			slog.Error("request failed", "err", err)
			return
		}

		comics, err := suggestComix.Run(
			req.language,
			req.query,
			req.limit,
		)

		if err != nil {
			err = Encode(w, r, http.StatusInternalServerError, err.Error())
			slog.Error("comix suggestion failed", "err", err)
			return
		}

		response := SuggestComixPicsResponse{
			Pics: make([]string, len(comics)),
		}

		for i, c := range comics {
			response.Pics[i] = c.URL
		}

		err = Encode(w, r, http.StatusOK, response)
		if err != nil {
			slog.Error("response failed", "err", err)
			return
		}
	}
}

type UpdateDatabaseAndIndexResponse struct {
	TotalComics int `json:"total_comics"`
	NewAdded    int `json:"new_added"`
}

func UpdateDatabaseAndIndexHandler(
	ctx context.Context,
	downloadComics *downloadComics.UseCase,
	buildIndex *buildIndex.UseCase,
	countComics *countComics.UseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countBeforeUpdate, err := countComics.Run()
		if err != nil {
			slog.Error("failed to count", "err", err)
			return
		}

		err = downloadComics.Run(ctx)
		if err != nil {
			slog.Error("comix download failed", "err", err)
			return
		}
		err = buildIndex.Run(ctx)
		if err != nil {
			slog.Error("building index failed", "err", err)
			return
		}

		countAfterUpdate, err := countComics.Run()
		if err != nil {
			slog.Error("failed to count", "err", err)
			return
		}

		response := UpdateDatabaseAndIndexResponse{
			TotalComics: countAfterUpdate,
			NewAdded:    countAfterUpdate - countBeforeUpdate,
		}

		slog.Info("total", "total", response.TotalComics)

		err = Encode(w, r, http.StatusOK, response)
		if err != nil {
			slog.Error("response failed", "err", err)
			return
		}
	}
}
