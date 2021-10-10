package app

import (
	"context"
	"github.com/pkg/errors"
	"github.com/torwig/promotions/internal/promotions/repo"
	"github.com/torwig/promotions/internal/promotions/repo/mysql"
	"github.com/torwig/promotions/internal/promotions/types"
)

var ErrPromotionNotFound = errors.New("promotion not found")

type Application struct {
	repo repo.Repository
}

func NewApplication(ctx context.Context) Application {
	rp, err := mysql.NewRepository()
	if err != nil {
		panic(err)
	}

	return Application{repo: rp}
}

func (a Application) GetPromotion(ctx context.Context, id string) (*types.Promotion, error) {
	p, err := a.repo.Get(ctx, id)
	if errors.Is(err, repo.ErrPromotionNotFound) {
		return nil, ErrPromotionNotFound
	} else if err != nil {
		return nil, err
	}

	return p, nil
}
