package http

import (
	"errors"
	"net/http"
	"strconv"

	"log/slog"

	countComics "github.com/toadharvard/goxkcd/internal/usecase/countcomics"
	downloadComics "github.com/toadharvard/goxkcd/internal/usecase/downloadcomics"
	suggestComix "github.com/toadharvard/goxkcd/internal/usecase/suggestcomix"
	"github.com/toadharvard/goxkcd/pkg/iso6391"
)

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = Encode(w, r, http.StatusOK, "pong")
	}
}

type SuggestComixPicsResponse struct {
	Pics []string `json:"pics"`
}

type SuggestComixPicsRequest struct {
	language iso6391.ISOCode6391
	search   string
	limit    int
}

var ErrMissingSearch = errors.New("missing search query")
var ErrInvalidLanguageCode = errors.New("invalid language code query")

func NewSuggestComixPicsRequest(r *http.Request) (req SuggestComixPicsRequest, err error) {
	languageStr := r.URL.Query().Get("language")

	if languageStr == "" {
		languageStr = "en"
	}

	language, err := iso6391.NewLanguage(
		languageStr,
	)

	if err != nil {
		return SuggestComixPicsRequest{}, ErrInvalidLanguageCode
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		return SuggestComixPicsRequest{}, ErrMissingSearch
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))

	if err != nil {
		limit = 10
	}

	return SuggestComixPicsRequest{
		language: language,
		search:   search,
		limit:    limit,
	}, nil
}

func SuggestComixPicsHandler(suggestComix *suggestComix.UseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := NewSuggestComixPicsRequest(r)
		if err != nil {
			slog.Error("request failed", "err", err)
			_ = Encode(w, r, http.StatusBadRequest, err.Error())
			return
		}

		comics, err := suggestComix.Run(
			req.language,
			req.search,
			req.limit,
		)

		if err != nil {
			slog.Error("comix suggestion failed", "err", err)
			w.WriteHeader(http.StatusOK)
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
	downloadComics *downloadComics.UseCase,
	countComics *countComics.UseCase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countBeforeUpdate, err := countComics.Run()
		if err != nil {
			slog.Error("failed to count", "err", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		err = downloadComics.Run(r.Context())
		if err != nil {
			slog.Error("comix download failed", "err", err)
			w.WriteHeader(http.StatusOK)
			return
		}

		countAfterUpdate, err := countComics.Run()
		if err != nil {
			slog.Error("failed to count", "err", err)
			w.WriteHeader(http.StatusOK)
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
