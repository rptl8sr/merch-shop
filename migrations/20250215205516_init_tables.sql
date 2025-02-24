-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id serial primary key,
    username varchar(255) unique not null,
    password_hash varchar(255) not null,
    coins_balance int not null default 1000 check (coins_balance >= 0),
    created_at timestamp default current_timestamp
);

create index if not exists idx_users_username on users(username);

create table if not exists transactions (
    id serial primary key,
    sender_id int references users(id),
    receiver_id int not null references users(id),
    amount int not null check (amount > 0),
    created_at timestamp default current_timestamp
);

create index if not exists idx_transactions_sender_id on transactions(sender_id, created_at);
create index if not exists idx_transactions_receiver_id on transactions(receiver_id, created_at);

create table if not exists merch_items (
    id serial primary key,
    item_name varchar(255) unique not null,
    price int not null check (price > 0),
    created_at timestamp default current_timestamp
);

create index if not exists idx_merch_items_name on merch_items(item_name);

create table if not exists purchases (
    id serial primary key,
    user_id int not null references users(id),
    item_id int not null references merch_items(id),
    quantity int not null check (quantity > 0),
    total_cost int not null check (total_cost > 0),
    created_at timestamp default current_timestamp
);

create index if not exists idx_purchases_user on purchases(user_id, created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists idx_purchases_user;
drop table if exists purchases;

drop index if exists idx_merch_items_name;
drop table if exists merch_items;

drop index if exists idx_transactions_sender_id;
drop index if exists idx_transactions_receiver_id;
drop table if exists transactions;

drop index if exists idx_users_username;
drop table if exists users;
-- +goose StatementEnd
