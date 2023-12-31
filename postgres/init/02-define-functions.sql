create or replace function l0.save_order(
	 _order_uid          uuid,
     _track_number       text,
     _entry              text,
     _d                  l0.delivery,
     _p                  l0.payment,
     _locale             l0.locale_value,
     _internal_signature text,
     _customer_id        text,
     _delivery_service   text,
     _shard_key          text,
     _sm_id              int,
     _date_created       timestamp,
     _oof_shard          text,
     _items              l0.item[]
  ) returns bool 
   language plpgsql
  as
$$
declare 
  _item l0.item;
  _rids text[];
begin
    with new_items as (
    insert into l0.item 
		(rid, chrt_id, nm_id, "name", brand, "size", track_number,
			price, sale, total_price, status)
		(select * from unnest(_items)) 
	returning rid
	),
	new_delivery as (
		insert into l0.delivery 
			("name", phone, zip, city, address, region, email)
		values
			(_d."name", _d.phone, _d.zip, _d.city, _d.address, _d.region, _d.email) 
		returning delivery_id
	),
	new_payment as (
		insert into l0.payment
			("transaction", request_id, currency, provider, amount, payment_dt, 
				bank, delivery_cost, goods_total, custom_fee)
		values
			(_p.transaction, _p.request_id, _p."currency", _p.provider,
				_p.amount, _p.payment_dt, _p.bank, _p.delivery_cost,
				_p.goods_total, _p.custom_fee)
		 returning payment_id
	),
	new_order as (
		insert into l0.orders 
			(order_uid, track_number, entry, delivery_id, payment_id,
				locale, internal_signature,  customer_id, delivery_service,
				shard_key, sm_id, date_created, oof_shard)
		values
			(_order_uid, _track_number, _entry,
				(select delivery_id from new_delivery),
				(select payment_id from new_payment),
				_locale, _internal_signature, _customer_id, _delivery_service,
				_shard_key, _sm_id, _date_created, _oof_shard)
	)
	insert into l0.orders_items (order_uid, item_rid)
		select _order_uid, i.rid
			from new_items i;
	return true;
end;
$$;

-- example usage
/*
select l0.save_order(
	'20354d7a-e4fe-47af-8ff6-187bca92f3f9'::uuid, 
	'WBILMTESTTRACK', 
	'WBIL',
	-- delivery
	('00000000-0000-0000-0000-000000000000', 'Test Testov', '+79720000000', '269809', 'Kiryat Mozkin', 'Ploshad Mira 15', 
		'Kraiot', 'test@gmail.com')::l0.delivery,
	-- payment
	('00000000-0000-0000-0000-000000000000', 'b563feb7b2b84b6test', '', 'USD'::currency, 'wbpay', 1817, '2021-07-21'::timestamp, 
		'alpha', 1500, 317, 0)::l0.payment,
	-- everything else
	'en'::l0.locale_value, 
	'', 
	'00000000-0000-0000-1000-000000000000', 
	'service_1',
	'9', 
	99, 
	'2023-03-01 06:22:19'::timestamp, 
	'1',
	-- items
	array[
		('ab4219087a764ae0btest', 9934930, 2389212, 'Mascaras', 'Vivienne Sabo', '0', 
			'WBILMTESTTRACK', 453, 30, 317, 202)
	]::l0.item[]
);
*/


create or replace function l0.get_all_orders() 
   returns table (order_data json)
   language plpgsql
  as
$$
begin
	return query
	with order_to_item_data as (
		select oi.order_uid, jsonb_agg(to_jsonb(i)) items
		from l0.orders_items oi
		join l0.item i on oi.item_rid = i.rid
		group by oi.order_uid
	),
	prepared_data as (
		select
			o.order_uid,
			o.track_number,
			o.entry,
			otid.items,
			to_jsonb(d) - 'delivery_id' delivery,
			(to_jsonb(p) - 'payment_id') || jsonb_build_object('payment_dt', extract(epoch from p.payment_dt)::int) payment,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shard_key,
			o.sm_id,
			to_char(o.date_created, 'YYYY-MM-DD"T"HH24:MI:SSOF'),
			o.oof_shard
		from l0.orders o
		join order_to_item_data otid on o.order_uid = otid.order_uid
		join l0.delivery d on d.delivery_id = o.delivery_id
		join l0.payment p on p.payment_id = o.payment_id
	)
	select row_to_json(d) as order_data
	from prepared_data d;
end;
$$;