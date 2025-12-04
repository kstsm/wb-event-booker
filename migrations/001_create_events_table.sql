-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    id                            UUID PRIMARY KEY,
    name                          VARCHAR(64) NOT NULL,
    date                          TIMESTAMPTZ NOT NULL,
    total_seats                   INT         NOT NULL CHECK (total_seats > 0),
    reserved_seats                INT         NOT NULL DEFAULT 0 CHECK (reserved_seats >= 0),
    booked_seats                  INT         NOT NULL DEFAULT 0 CHECK (booked_seats >= 0),
    booking_lifetime              INT         NOT NULL CHECK (booking_lifetime > 0),
    requires_payment_confirmation BOOLEAN     NOT NULL DEFAULT true,
    created_at                    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT seats_check CHECK (reserved_seats + booked_seats <= total_seats)
);

CREATE INDEX IF NOT EXISTS idx_events_date ON events (date);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events (created_at);
CREATE INDEX IF NOT EXISTS idx_events_name ON events (name);

-- +goose Down
DROP TABLE IF EXISTS events;




