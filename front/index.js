const LayoutSchema = {
  delivery: {
    name:   'Name', 
    phone:  'Phone',
    zip:    'Zip-code',
    city:   'City',
    region: 'Region',
    email:  'Email'
  },
  payment: {
    amount:        'Amount',
    bank:          'Bank',
    currency:      'Currency',
    custom_fee:    'Custom fee',
    delivery_cost: 'Delivery cost',
    goods_total:   'Goods Total',
    payment_dt:    'Payment time',
    provider:      'Provider',
    request_id:    'Request ID',
    transaction:   'Transaction'
  },
  order: {
    customer_id:        'Customer ID',
    date_created:       'Creation Date',
    delivery_service:   'Delivery Service',
    entry:              'Entry',
    internal_signature: 'Internal Signature',
    locale:             'Locale',
    oof_shard:          'Shard OOF',
    order_uid:          'Order UID',
    shardkey:           'Shard Key',
    sm_id:              'SM ID',
    track_number:       'Track Number'
  },
  item: {
    brand: 'Brand',
    chrt_id: 'Vendor code',
    name: 'Name',
    nm_id: 'NM ID',
    price: 'Price',
    rid: 'Item ID',
    sale: 'Sale',
    size: 'Size',
    status: 'Status',
    total_price: 'Total Price',
    track_number: 'Track Number'
  }
}