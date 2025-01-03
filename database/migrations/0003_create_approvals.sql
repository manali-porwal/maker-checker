-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS 
    message_approvals (
                message_id BIGINT not null,
                approver_id BIGINT not null,
                approved BOOLEAN not null,
                comment TEXT,
                created_at TIMESTAMP NOT NULL,
                primary key(message_id, approver_id),
                constraint fk_messages foreign key(message_id) references messages(id) on delete cascade,
                constraint fk_users foreign key(approver_id) references users(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists message_approvals;
-- +goose StatementEnd