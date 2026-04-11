package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChernykhITMO/Wishlist-API/internal/booking/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type bookingRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *bookingRepository {
	return &bookingRepository{
		pool: pool,
	}
}

func (r *bookingRepository) Create(ctx context.Context, token, giftID, bookingID uuid.UUID) error {
	const query = `
	INSERT INTO bookings (id, gift_id)
	SELECT $3, g.id
	FROM gifts AS g
	JOIN wishlists AS w ON w.id = g.wishlist_id
	WHERE w.token = $1 AND g.id = $2
	`

	tag, err := r.pool.Exec(ctx, query, token, giftID, bookingID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return domain.ErrBookingAlreadyExists
			}
		}
		return fmt.Errorf("create booking: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrGiftNotFound
	}
	return nil
}
