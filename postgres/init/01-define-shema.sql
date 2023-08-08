create schema L0;

create type L0.locale_value as enum ('ru', 'en');
create type L0.currency as enum ('RUB', 'USD', 'EUR', 'CNY');


create table L0.delivery (
	delivery_id uuid primary key default(gen_random_uuid()),
	name        text not null,
	phone       text not null check(phone ~ '^\+\d{11}$'),
	zip         text not null check(zip ~ '^\d{6}$'), -- assuming russian zip format
	city        text not null check(length(city) between 1 and 500),
	address     text not null check(length(address) between 1 and 500),
	region      text not null check(length(region) between 1 and 100),
	email       text not null check(email ~ '^[\w\.-]+@([\w-]+\.)+\w+$')
);

create table L0.item (
	rid     text primary key, -- specific item code
	chrt_id int not null, -- vendor code
	nm_id   int not null, -- WB id
	
	name  text not null,
	brand text not null,
	size  text not null,
	
	track_number text not null,
	
	price       int not null check(price > 0),
	sale        int not null check(sale between 0 and 100),
	total_price int not null check(total_price = price * (100 - sale) / 100),
	
	status int not null
);

create table L0.payment (
	payment_id uuid primary key default(gen_random_uuid()),
	transaction text not null,
	request_id text,
	currency L0.currency not null,
	provider text not null,
	amount int not null check (amount > 0),
	payment_dt timestamp not null,
	bank text not null,
	delivery_cost int not null check (delivery_cost >= 0),
	goods_total int not null check (goods_total >= 0),
	custom_fee int not null check (custom_fee >= 0)
);

create table L0.orders (
	order_uid uuid primary key default(gen_random_uuid()),
	track_number text unique not null,
	entry text not null, -- what is this?
	delivery_id uuid not null,
	payment_id uuid,
    locale L0.locale_value not null,
    internal_signature text,
    customer_id text,
    delivery_service text not null,
    shard_key text not null,
    sm_id int not null,
    date_created timestamp not null,
    oof_shard text not null,
    
    foreign key (delivery_id) references L0.delivery(delivery_id),
    foreign key (payment_id) references L0.payment(payment_id)
);


create table L0.orders_items (
	order_uid uuid,
	item_rid text,

	foreign key (order_uid) references L0.orders(order_uid),
	foreign key (item_rid) references L0.item(rid),
	primary key (order_uid, item_rid)
);

create index items_by_order on L0.orders_items(order_uid);
