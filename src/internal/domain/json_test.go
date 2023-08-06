package domain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

func compressJsonString(jsonString string) string {
	lines := strings.Split(jsonString, "\n")
	for i, line := range lines {
		line = strings.Trim(line, " \t")
		line = strings.ReplaceAll(line, `": `, `":`)
		lines[i] = line
	}
	return strings.Join(lines, "")
}

func testUnmarshal[T any](t *testing.T, data string, expected T) {
	var result T
	json.Unmarshal([]byte(data), &result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("error unmarshalling %T: expected %v, got %v", result, expected, result)
	}
}

func testMarshal[T any](t *testing.T, data T, expected string) {
	expected = compressJsonString(expected)
	result, err := json.Marshal(data)
	if err != nil {
		err = fmt.Errorf("error marshalling %T: %w", data, err)
		t.Error(err)
	}
	if !reflect.DeepEqual(string(result), expected) {
		t.Errorf("error marshalling %T: expected %v, got %v", data, expected, string(result))
	}
}

func testJSON[T any](t *testing.T, object T, jsonString string) {
	testMarshal[T](t, object, jsonString)
	testUnmarshal[T](t, jsonString, object)
}

var expectedOrder = Order{
	OrderUid:    "b563feb7b2b84b6test",
	TrackNumber: "WBILMTESTTRACK",
	Entry:       "WBIL",
	Delivery: Delivery{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Kiryat Mozkin",
		Address: "Ploshad Mira 15",
		Region:  "Kraiot",
		Email:   "test@gmail.com",
	},
	Payment: Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestId:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    Timestamp{time.Unix(1637907727, 0)},
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	},
	Locale:            "en",
	InternalSignature: "",
	CustomerId:        "test",
	DeliveryService:   "meest",
	Shard_key:         "9",
	SmId:              99,
	DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
	Oof_shard:         "1",
	Items: []Item{{
		RId:         "ab4219087a764ae0btest",
		ChrtId:      9934930,
		NmId:        2389212,
		Name:        "Mascaras",
		Brand:       "Vivienne Sabo",
		Size:        "0",
		TrackNumber: "WBILMTESTTRACK",
		Price:       453,
		Sale:        30,
		TotalPrice:  317,
		Status:      202,
	}},
}

func TestPayment(t *testing.T) {
	data := `{
	  "transaction": "b563feb7b2b84b6test",
	  "request_id": "",
	  "currency": "USD",
	  "provider": "wbpay",
	  "amount": 1817,
	  "payment_dt": 1637907727,
	  "bank": "alpha",
	  "delivery_cost": 1500,
	  "goods_total": 317,
	  "custom_fee": 0
	}`
	expected := expectedOrder.Payment
	testJSON[Payment](t, expected, data)
}

func TestItem(t *testing.T) {
	data := `{
      "chrt_id": 9934930,
      "track_number": "WBILMTESTTRACK",
      "price": 453,
      "rid": "ab4219087a764ae0btest",
      "name": "Mascaras",
      "sale": 30,
      "size": "0",
      "total_price": 317,
      "nm_id": 2389212,
      "brand": "Vivienne Sabo",
      "status": 202
    }`
	expected := expectedOrder.Items[0]
	testJSON[Item](t, expected, data)
}

func TestDeliveryUnmarshalling(t *testing.T) {
	data := `{
	  "name": "Test Testov",
	  "phone": "+9720000000",
	  "zip": "2639809",
	  "city": "Kiryat Mozkin",
	  "address": "Ploshad Mira 15",
	  "region": "Kraiot",
	  "email": "test@gmail.com"
	}`
	expected := expectedOrder.Delivery
	testJSON[Delivery](t, expected, data)
}

func TestOrderUnmarshalling(t *testing.T) {
	data := `
	{
	  "order_uid": "b563feb7b2b84b6test",
	  "track_number": "WBILMTESTTRACK",
	  "entry": "WBIL",
	  "delivery": {
	    "name": "Test Testov",
	    "phone": "+9720000000",
	    "zip": "2639809",
	    "city": "Kiryat Mozkin",
	    "address": "Ploshad Mira 15",
	    "region": "Kraiot",
	    "email": "test@gmail.com"
	  },
	  "payment": {
	    "transaction": "b563feb7b2b84b6test",
	    "request_id": "",
	    "currency": "USD",
	    "provider": "wbpay",
	    "amount": 1817,
	    "payment_dt": 1637907727,
	    "bank": "alpha",
	    "delivery_cost": 1500,
	    "goods_total": 317,
	    "custom_fee": 0
	  },
	  "items": [
	    {
	      "chrt_id": 9934930,
	      "track_number": "WBILMTESTTRACK",
	      "price": 453,
	      "rid": "ab4219087a764ae0btest",
	      "name": "Mascaras",
	      "sale": 30,
	      "size": "0",
	      "total_price": 317,
	      "nm_id": 2389212,
	      "brand": "Vivienne Sabo",
	      "status": 202
	    }
	  ],
	  "locale": "en",
	  "internal_signature": "",
	  "customer_id": "test",
	  "delivery_service": "meest",
	  "shardkey": "9",
	  "sm_id": 99,
	  "date_created": "2021-11-26T06:22:19Z",
	  "oof_shard": "1"
	}`
	expected := expectedOrder
	testJSON[Order](t, expected, data)
}
