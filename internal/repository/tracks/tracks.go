package tracks

import (
	"time"

	"github.com/antongoncharik/crypto-knight-api/internal/entity/track"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Tracks struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Tracks {
	return &Tracks{db: db}
}

func (t *Tracks) GetAll(queryParams track.QueryParams) ([]track.Track, error) {
	tracksData := []track.Track{}

	if queryParams.Full && queryParams.Symbol != "" {
		err := t.db.Select(&tracksData, "select symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, is_order, high_created_at, low_created_at, high_prices, low_prices from tracks where created_at between $1 and $2 and symbol = $3 order by created_at desc", queryParams.From, queryParams.To, queryParams.Symbol)

		return tracksData, err
	}

	if queryParams.Full && queryParams.Symbol == "" {
		err := t.db.Select(&tracksData, "select symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, is_order, high_created_at, low_created_at, high_prices, low_prices from tracks where created_at between $1 and $2 order by created_at desc", queryParams.From, queryParams.To)

		return tracksData, err
	}

	if queryParams.Symbol != "" {
		err := t.db.Select(&tracksData, "with ranked_tracks as (select id, symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, lag(high_price) over (partition by symbol order by created_at) as prev_high_price, lag(low_price) over (partition by symbol order by created_at) as prev_low_price, high_created_at, low_created_at, high_prices, low_prices from tracks where created_at between $1 and $2 and symbol = $3) select symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, high_created_at, low_created_at, high_prices, low_prices from ranked_tracks where high_price != prev_high_price or low_price != prev_low_price or prev_high_price is null order by created_at desc", queryParams.From, queryParams.To, queryParams.Symbol)

		return tracksData, err
	}

	err := t.db.Select(&tracksData, "with ranked_tracks as (select id, symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, lag(high_price) over (partition by symbol order by created_at) as prev_high_price, lag(low_price) over (partition by symbol order by created_at) as prev_low_price, high_created_at, low_created_at, high_prices, low_prices from tracks where created_at between $1 and $2) select symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, high_created_at, low_created_at, high_prices, low_prices from ranked_tracks where high_price != prev_high_price or low_price != prev_low_price or prev_high_price is null order by created_at desc", queryParams.From, queryParams.To)

	return tracksData, err
}

func (t *Tracks) Create(track track.Track) error {
	var err error
	if (track.CreatedAt == time.Time{}) {
		_, err = t.db.Exec(`INSERT INTO tracks (symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, is_order, high_created_at, low_created_at, high_prices, low_prices)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`, track.Symbol, track.HighPrice, track.LowPrice, track.HighPrice1, track.LowPrice1, track.HighPrice2, track.LowPrice2, track.HighPrice3, track.LowPrice3, track.IsOrder, track.HighCreatedAt, track.LowCreatedAt, track.HighPrices, track.LowPrices)
	} else {
		_, err = t.db.Exec(`INSERT INTO tracks (symbol, high_price, low_price, high_price_1, low_price_1, high_price_2, low_price_2, high_price_3, low_price_3, created_at, is_order, high_created_at, low_created_at, high_prices, low_prices)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, track.Symbol, track.HighPrice, track.LowPrice, track.HighPrice1, track.LowPrice1, track.HighPrice2, track.LowPrice2, track.HighPrice3, track.LowPrice3, track.CreatedAt, track.IsOrder, track.HighCreatedAt, track.LowCreatedAt, track.HighPrices, track.LowPrices)
	}

	return err
}
