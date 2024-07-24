create table if not exists account
(
    id         serial primary key,
    user_id    int       not null unique,
    balance    float     not null default 0,
    created_at timestamp not null default now()
);

create table if not exists reservation
(
    id         serial primary key,
    user_id    int       not null references account (user_id),
    product_id int       not null, -- not unique because user can buy same products
    order_id   int       not null, -- not unique because user can buy more than 1 product in order
    amount     float     not null,
    created_at timestamp not null default now()
);

create table if not exists operation
(
    id         serial primary key,
    user_id    int       not null references account (user_id),
    product_id int                default null,
    order_id   int                default null,
    amount     float     not null,
    type       varchar   not null,
    created_at timestamp not null default now()
);
