-- +goose Up
CREATE table user (
    id      serial primary key,
    login   text
);

INSERT INTO user (login, description, meta, updated_at)
VALUES
    ('admin');

CREATE table event (
    id              serial primary key,
	name            text not null,
	date_time       timestamptz not null,
	end_date_time   timestamptz not null,
	description     text,
	user_id         integer not null,
    CONSTRAINT      event_user_id_fk FOREIGN KEY(user_id) REFERENCES user(id),
	send_before     timestamptz
);

-- +goose Down
drop table event;
drop table user;
