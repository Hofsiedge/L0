package postgres

import (
	"context"
	"database/sql/driver"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/Hofsiedge/l0/internal/domain"
)

const uuidPlaceholder = "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"

type ItemRecord domain.Item

func (i ItemRecord) Value() (driver.Value, error) {
	record := fmt.Sprintf(
		`("%s",%d,%d,"%s","%s","%s","%s",%d,%d,%d,%d)`,
		i.RId, i.ChrtId, i.NmId, i.Name, i.Brand, i.Size,
		i.TrackNumber, i.Price, i.Sale, i.TotalPrice, i.Status)
	return record, nil
}

// &Orders implements repo.Repo[domain.Order, string]
type Orders struct {
	conn *pgx.Conn
}

func NewOrderRepo(dbURL string) (*Orders, error) {
	var (
		conn *pgx.Conn
		err  error
	)
	for attempts_left := 3; attempts_left > 0; attempts_left-- {
		conn, err = pgx.Connect(context.Background(), dbURL)
		if err == nil {
			break
		}
		err = fmt.Errorf("could not connect to postgres: %w", err)
		log.Printf("%v. trying again after 2 seconds\n", err)
		time.Sleep(2 * time.Second)
	}
	orders := Orders{conn}
	return &orders, err
}

func (o *Orders) Close() error {
	return o.conn.Close(context.Background())
}

func (o *Orders) Get(id string) (domain.Order, error) {
	return domain.Order{}, nil
}

func (o *Orders) List() ([]string, error) {
	return make([]string, 0), nil
}

func (o *Orders) GetAll() ([]domain.Order, error) {
	orders := make([]domain.Order, 0)

	rows, err := o.conn.Query(
		context.Background(),
		`select order_data from l0.get_all_orders();`)
	if err != nil {
		err = fmt.Errorf("error reading orders: %w", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order); err != nil {
			err = fmt.Errorf("error scanning an Order: %w", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	log.Printf("read %d orders from the database\n", len(orders))
	return orders, nil
}

// this is the ugliest code I've ever written
//
// I just hope I missed something and there is a way
// to handle this with driver.Valuer and pgx
func (o *Orders) Save(order domain.Order) error {
	deliveryArgs := []any{
		uuidPlaceholder,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	}
	paymentArgs := []any{
		uuidPlaceholder,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt.Format("'2006-1-2 15:4:5'"),
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	}

	args := make([]any, 0)
	args = append(args, order.OrderUid, order.TrackNumber, order.Entry)
	args = append(args, deliveryArgs...)
	args = append(args, paymentArgs...)
	leftoverArgs := []any{
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shard_key,
		order.SmId,
		order.DateCreated.Format("'2006-1-2 15:4:5'"),
		order.Oof_shard,
	}
	args = append(args, leftoverArgs...)
	positions := make([]string, len(order.Items))
	for i := range order.Items {
		positions[i] = fmt.Sprintf("$%d", i+31)
	}

	itemRecords := make([]any, len(order.Items))
	for i, item := range order.Items {
		itemRecords[i] = ItemRecord(item)
	}
	args = append(args, itemRecords...)

	// this is ugly, but pgx can't handle record types and
	// arrays of records with driver.Valuer :(
	sqlString := `select l0.save_order($1::uuid, $2, $3, 
			 ($4::uuid, $5, $6, $7, $8, $9, $10, $11)::l0.delivery,
			 ($12::uuid, $13, $14, $15::l0.currency, $16, $17, 
				$18, $19, $20, $21, $22)::l0.payment,
			 $23::l0.locale_value,
			 $24, $25, $26, $27, $28, $29::timestamp, $30,` +
		fmt.Sprintf("array[%s]::l0.item[]);", strings.Join(positions, ", "))
	_, err := o.conn.Exec(
		context.Background(),
		sqlString,
		args...)
	if err != nil {
		err = fmt.Errorf("SQL error saving an Order: %w", err)
		return err
	}
	return nil
}
