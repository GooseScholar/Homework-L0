package database

import (
	"context"
	"errors"
	"homework-l0/internal/models"
	"log"

	"github.com/jackc/pgx"
)

//Запись новых даных в бд
func (db *DB) PutOrder(ctx context.Context, ord *models.Orders) (err error) {
	const query = `
		BEGIN TRY
		BEGIN TRANSACTION
  			INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature,
    			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
  			VALUES ($1, $2, $3, );

  		INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email) 
		VALUES ();

  		INSERT INTO payment(order_uid, transaction, request_id, currency, provider, amount,
    		payment_dt, bank, delivery_cost, goods_total, custom_fee)
  		VALUES ();

  		INSERT INTO items(order_uid, chrt_id, track_number, price, rid, name, sale, size,
    		total_price, nm_id, brand, status)
  		VALUES ();

			END TRY
      			BEGIN CATCH 
        		ROLLBACK TRANSACTION
        		SELECT ERROR_NUMBER(), ERROR_MESSAGE
				RETURN
      		END CATCH
		COMMIT TRANSACTION
`

	return
}

//Получение новых данных из бд
func (db *DB) GetOrder(ctx context.Context, order_uid string) (ord *models.Orders, err error) {
	const (
		getOrders = `
		SELECT * FROM orders o
		WHERE o.order_uid = &1
		`

		getDelivery = `
		SELECT name, phone, zip, city, address, region, email FROM delivery d
		WHERE o.order_uid = &1
		`

		getPayment = `
		SELECT * FROM payment p
		WHERE o.order_uid = &1
		`

		getItems = `
		SELECT * FROM items i
		WHERE o.order_uid = &1
		`
	)

	errOrders := db.pool.QueryRow(ctx, getOrders, order_uid).Scan(
		&ord.Order_uid,
		&ord.Track_number,
		&ord.Entry,
		&ord.Locale,
		&ord.Internal_signature,
		&ord.Customer_id,
		&ord.Delivery_service,
		&ord.Shardkey,
		&ord.Sm_id,
		&ord.Date_created,
		&ord.Oof_shard,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		errOrders = nil
	}
	errDelivery := db.pool.QueryRow(ctx, getDelivery, order_uid).Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		errDelivery = nil
	} else {
		return
	}
	errPayment := db.pool.QueryRow(ctx, getPayment, order_uid).Scan(
		&ord.Payment.Transaction,
		&ord.Payment.Request_id,
		&ord.Payment.Currency,
		&ord.Payment.Provider,
		&ord.Payment.Amount,
		&ord.Payment.Payment_dt,
		&ord.Payment.Bank,
		&ord.Payment.Delivery_cost,
		&ord.Payment.Goods_total,
		&ord.Payment.Custom_fee,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		errPayment = nil
	} else {
		return
	}
	var item models.Item

	var errItem = db.pool.QueryRow(ctx, getItems, order_uid).Scan(
		item.Chrt_id,
		item.Track_number,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.Total_price,
		item.Nm_id,
		item.Brand,
		item.Status,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		errItem = nil
	}

	if (errOrders != nil) || (errDelivery != nil) || (errPayment != nil) || (errItem != nil) {
		log.Printf("table orders <%v> | table delivery <%v> | table payment <%v> | table items <%v>",
			errOrders, errDelivery, errPayment, errItem)
	}

	ord.Items = append(ord.Items, item)

	return

}
