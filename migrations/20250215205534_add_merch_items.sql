-- +goose Up
-- +goose StatementBegin
insert into merch_items (item_name, price) values
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500)
on conflict (item_name) do nothing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from merch_items where item_name in (
    't-shirt',
    'cup',
    'book',
    'pen',
    'powerbank',
    'hoody',
    'umbrella',
    'socks',
    'wallet',
    'pink-hoody'
);
-- +goose StatementEnd
