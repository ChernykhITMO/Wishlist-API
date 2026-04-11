package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ChernykhITMO/Wishlist-API/internal/gifts/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type giftRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *giftRepository {
	return &giftRepository{
		pool: pool,
	}
}

func (r *giftRepository) Create(ctx context.Context, gift domain.Gift) error {
	const query = `
	INSERT INTO gifts (id, wishlist_id, name, description, link, priority)
	SELECT $1, w.id, $4, $5, $6, $7
	FROM wishlists AS w
	WHERE w.id = $2 AND w.user_id = $3
	`

	tag, err := r.pool.Exec(ctx, query,
		gift.ID, gift.WishlistID,
		gift.UserID,
		gift.Name, gift.Description,
		gift.Link, gift.Priority,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return domain.ErrGiftAlreadyExists
			}
		}

		return fmt.Errorf("create gift: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrGiftNotFound
	}

	return nil
}

func (r *giftRepository) GetByID(ctx context.Context, giftID, wishlistID, userID uuid.UUID) (domain.Gift, error) {
	const query = `
	SELECT g.id, g.wishlist_id, g.name, g.description, g.link, g.priority
	FROM gifts AS g
	JOIN wishlists AS w ON w.id = g.wishlist_id
	WHERE g.id = $1 AND w.id = $2 AND w.user_id = $3
	`

	gift, err := scanGift(r.pool.QueryRow(ctx, query, giftID, wishlistID, userID))
	if err != nil {
		return domain.Gift{}, fmt.Errorf("get gift by id: %w", err)
	}

	return gift, nil
}

func (r *giftRepository) ListByWishlistID(ctx context.Context, wishlistID, userID uuid.UUID) ([]domain.Gift, error) {
	const query = `
	SELECT g.id, g.wishlist_id, g.name, g.description, g.link, g.priority
	FROM gifts AS g
	JOIN wishlists AS w ON w.id = g.wishlist_id
	WHERE w.id = $1 AND w.user_id = $2
	ORDER BY g.priority DESC, g.name ASC
	`

	rows, err := r.pool.Query(ctx, query, wishlistID, userID)
	if err != nil {
		return nil, fmt.Errorf("list gifts by wishlist id: %w", err)
	}
	defer rows.Close()

	gifts := make([]domain.Gift, 0)
	for rows.Next() {
		gift, err := scanGift(rows)
		if err != nil {
			return nil, fmt.Errorf("scan gift: %w", err)
		}
		gifts = append(gifts, gift)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gifts: %w", err)
	}

	return gifts, nil
}

func (r *giftRepository) Update(ctx context.Context, gift domain.Gift) error {
	const query = `
	UPDATE gifts AS g
	SET name = $4,description = $5, link = $6, priority = $7
	FROM wishlists AS w
	WHERE g.id = $1 AND g.wishlist_id = $2 AND w.id = g.wishlist_id AND w.user_id = $3
	`

	tag, err := r.pool.Exec(
		ctx,
		query,
		gift.ID,
		gift.WishlistID,
		gift.UserID,
		gift.Name,
		gift.Description,
		gift.Link,
		gift.Priority,
	)
	if err != nil {
		return fmt.Errorf("update gift: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrGiftNotFound
	}

	return nil
}

func (r *giftRepository) Delete(ctx context.Context, giftID, wishlistID, userID uuid.UUID) error {
	const query = `
	DELETE FROM gifts AS g
	USING wishlists AS w
	WHERE g.id = $1
	  AND g.wishlist_id = $2
	  AND w.id = g.wishlist_id
	  AND w.user_id = $3
	`

	tag, err := r.pool.Exec(ctx, query, giftID, wishlistID, userID)
	if err != nil {
		return fmt.Errorf("delete gift: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrGiftNotFound
	}

	return nil
}

func (r *giftRepository) ListByToken(ctx context.Context, token uuid.UUID) ([]domain.Gift, error) {
	const query = `
	SELECT g.id, g.wishlist_id, g.name, g.description, g.link, g.priority 
	FROM gifts AS g
	JOIN wishlists AS w ON w.id = g.wishlist_id
	WHERE w.token = $1`

	rows, err := r.pool.Query(ctx, query, token)
	if err != nil {
		return nil, fmt.Errorf("list gifts by wishlist id: %w", err)
	}
	defer rows.Close()

	gifts := make([]domain.Gift, 0)
	for rows.Next() {
		gift, err := scanGift(rows)
		if err != nil {
			return nil, fmt.Errorf("scan gift: %w", err)
		}
		gifts = append(gifts, gift)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gifts: %w", err)
	}

	return gifts, nil
}

type giftScanner interface {
	Scan(dest ...any) error
}

func scanGift(scanner giftScanner) (domain.Gift, error) {
	var gift domain.Gift

	err := scanner.Scan(
		&gift.ID,
		&gift.WishlistID,
		&gift.Name,
		&gift.Description,
		&gift.Link,
		&gift.Priority,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Gift{}, domain.ErrGiftNotFound
		}
		return domain.Gift{}, err
	}

	return gift, nil
}
