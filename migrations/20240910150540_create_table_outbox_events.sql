-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox_events (
                               id UUID PRIMARY KEY,
                               aggregate_type VARCHAR(255) NOT NULL,
                               aggregate_id UUID NOT NULL,
                               event_type VARCHAR(255) NOT NULL,
                               payload JSONB NOT NULL,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_events_created_at ON outbox_events(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox_events;
-- +goose StatementEnd
