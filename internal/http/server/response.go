package server

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/torwig/promotions/internal/promotions/app"
	"github.com/torwig/promotions/internal/promotions/types"
)

type Promotion struct {
	ID             string  `json:"id"`
	Price          float64 `json:"price"`
	ExpirationDate string  `json:"expiration_date"`
}

func promotionResponse(p types.Promotion) Promotion {
	return Promotion{
		ID:             p.ID,
		Price:          p.Price,
		ExpirationDate: p.ExpirationDate.Format("2006-01-02 15:04:05"),
	}
}

func respond(w http.ResponseWriter, r *http.Request, data interface{}) {
	bb, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bb)
}

func respondWithError(err error, w http.ResponseWriter, r *http.Request) {
	if errors.Is(err, app.ErrPromotionNotFound) {
		http.NotFound(w, r)
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
