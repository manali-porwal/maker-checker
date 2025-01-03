-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
        id BIGSERIAL PRIMARY KEY,
        content TEXT NOT NULL,
        recipient VARCHAR(255) NOT NULL,
        status VARCHAR(50) NOT NULL,
        required_approvals INT NOT NULL,
        approval_count INT,
        rejection_count INT,
        created_by BIGINT REFERENCES users(id),
        created_at TIMESTAMP NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists messages;
-- +goose StatementEnd