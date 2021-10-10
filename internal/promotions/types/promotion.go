package types

import "time"

type Promotion struct {
	ID             string    `db:"id"`
	Price          float64   `db:"price"`
	ExpirationDate time.Time `db:"expiration_date"`
}
