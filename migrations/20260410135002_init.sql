-- +goose Up
CREATE TABLE users
(
    id            uuid PRIMARY KEY,
    email         text NOT NULL UNIQUE,
    password_hash text NOT NULL
);

CREATE TABLE wishlists
(
    id          uuid PRIMARY KEY,
    user_id     uuid        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token       uuid UNIQUE NOT NULL,
    name_event  text        NOT NULL,
    description text        NOT NULL,
    date_event  timestamptz NOT NULL
);

CREATE TABLE gifts
(
    id          uuid PRIMARY KEY,
    wishlist_id uuid NOT NULL REFERENCES wishlists (id) ON DELETE CASCADE,
    name        text NOT NULL,
    description text NOT NULL,
    link        text NOT NULL,
    priority    int   NOT NULL CHECK (priority >= 1 AND priority <= 5)
);

CREATE TABLE bookings
(
    id         uuid PRIMARY KEY,
    gift_id    uuid        NOT NULL UNIQUE REFERENCES gifts (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS gifts;
DROP TABLE IF EXISTS wishlists;
DROP TABLE IF EXISTS users;