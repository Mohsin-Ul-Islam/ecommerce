create table if not exists customers (
                                         id uuid primary key not null default gen_random_uuid(),
                                         email varchar unique not null,
                                         phone_number varchar unique not null,
                                         first_name varchar not null,
                                         last_name varchar not null,
                                         created_at timestamp with time zone not null default current_timestamp,
                                         updated_at timestamp with time zone not null default current_timestamp
);

create table if not exists products (
                                        id uuid primary key not null default gen_random_uuid(),
                                        title varchar not null,
                                        price double precision not null,
                                        available_quantity bigint not null,
                                        image_url varchar not null,
                                        description varchar not null,
                                        created_at timestamp with time zone not null default current_timestamp,
                                        updated_at timestamp with time zone not null default current_timestamp
);

create table if not exists transactions (
                                            id uuid primary key not null default gen_random_uuid(),
                                            customer_id uuid not null references customers(id),
                                            product_id uuid not null references products(id),
                                            quantity bigint not null,
                                            total_price double precision not null,
                                            created_at timestamp with time zone not null default current_timestamp
);

create table if not exists wallets (
                                            id uuid primary key not null default gen_random_uuid(),
                                            customer_id uuid not null references customers(id),
                                            balance double not null,
                                            created_at timestamp with time zone not null default current_timestamp,
                                            updated_at timestamp with time zone not null default current_timestamp
);

