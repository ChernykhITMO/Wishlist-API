package repository

import (
	"context"
	"errors"
	"github.com/ChernykhITMO/Wishlist-API/internal/wishlists/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type wishlistRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *wishlistRepository {
	return &wishlistRepository{
		pool: pool,
	}
}

func (r *wishlistRepository) Create(ctx context.Context, wishlist domain.Wishlist) error {
	const query = `
	INSERT INTO wishlists (
		id, user_id, token, name_event, description, date_event
   	) 
	VALUES ($1, $2, $3, $4, $5, $6) `

	_, err := r.pool.Exec(
		ctx, query,
		wishlist.ID,
		wishlist.UserID,
		wishlist.Token,
		wishlist.NameEvent,
		wishlist.Description,
		wishlist.DateEvent,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return domain.ErrWishlistAlreadyExists
			}
		}
		return err
	}

	return nil
}

func (r *wishlistRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (domain.Wishlist, error) {
	const query = `
	SELECT id, user_id, token, name_event, description, date_event
	FROM wishlists
	WHERE id = $1 AND user_id = $2
	`

	wishlist, err := scanWishlist(r.pool.QueryRow(ctx, query, id, userID))
	if err != nil {
		return domain.Wishlist{}, err
	}

	return wishlist, nil
}

func (r *wishlistRepository) GetByToken(ctx context.Context, token uuid.UUID) (domain.Wishlist, error) {
	const query = `
	SELECT id, user_id, token, name_event, description, date_event
	FROM wishlists
	WHERE token = $1
	`

	wishlist, err := scanWishlist(r.pool.QueryRow(ctx, query, token))
	if err != nil {
		return domain.Wishlist{}, err
	}

	return wishlist, nil
}

func (r *wishlistRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	const query = `
	SELECT id, user_id, token, name_event, description, date_event
	FROM wishlists
	WHERE user_id = $1
	ORDER BY date_event DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wishlists := make([]domain.Wishlist, 0)
	for rows.Next() {
		wishlist, err := scanWishlist(rows)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wishlists, nil
}

func (r *wishlistRepository) Update(ctx context.Context, wishlist domain.Wishlist) error {
	const query = `
	UPDATE wishlists
	SET name_event = $3,
	    description = $4,
	    date_event = $5
	WHERE id = $1 AND user_id = $2
	`

	tag, err := r.pool.Exec(
		ctx,
		query,
		wishlist.ID,
		wishlist.UserID,
		wishlist.NameEvent,
		wishlist.Description,
		wishlist.DateEvent,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrWishlistNotFound
	}

	return nil
}

func (r *wishlistRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	const query = `
	DELETE FROM wishlists
	WHERE id = $1 AND user_id = $2
	`

	tag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrWishlistNotFound
	}

	return nil
}

type wishlistScanner interface {
	Scan(dest ...any) error
}

func scanWishlist(scanner wishlistScanner) (domain.Wishlist, error) {
	var wishlist domain.Wishlist

	err := scanner.Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.Token,
		&wishlist.NameEvent,
		&wishlist.Description,
		&wishlist.DateEvent,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wishlist{}, domain.ErrWishlistNotFound
		}
		return domain.Wishlist{}, err
	}

	return wishlist, nil
}
