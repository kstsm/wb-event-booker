-- +goose Up

CREATE TYPE booking_status AS ENUM ('reserved', 'confirmed', 'cancelled');

CREATE TABLE IF NOT EXISTS bookings
(
    id         UUID PRIMARY KEY,
    event_id   UUID           NOT NULL REFERENCES events (id) ON DELETE CASCADE,
    user_id    UUID           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status     booking_status NOT NULL,
    deadline   TIMESTAMPTZ    NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_event_id ON bookings (event_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings (status);
CREATE INDEX IF NOT EXISTS idx_bookings_deadline ON bookings (deadline);
CREATE INDEX IF NOT EXISTS idx_bookings_event_status ON bookings (event_id, status);

-- +goose Down
DROP TABLE IF EXISTS bookings;
DROP TYPE IF EXISTS booking_status;





