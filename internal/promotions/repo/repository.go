package repo

import (
	"context"
	"errors"
	"github.com/torwig/promotions/internal/promotions/types"
)

var ErrPromotionNotFound = errors.New("promotion not found")

type Repository interface {
	Get(ctx context.Context, id string) (*types.Promotion, error)
}
