package database

import (
	"context"
	"homework-l0/internal/models"
	"log"
	"time"

	"homework-l0/internal/cache"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v4"
)

//Запись новых даных в бд
func (db *DB) PutOrder(ctx context.Context, ord *models.Orders) (err error) {

	//tx, err := db.pool.Begin(ctx)
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})

	log.Printf("tx %v\n", tx)

	//b := &pgx.Batch{}

	query := `
		INSERT INTO orders
			(order_uid, track_number, entry, locale, internal_signature,customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard, time_of_creation)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`
	_, err = tx.Exec(ctx,
		query,
		ord.Order_uid,
		ord.Track_number,
		ord.Entry,
		ord.Locale,
		ord.Internal_signature,
		ord.Customer_id,
		ord.Delivery_service,
		ord.Shardkey,
		ord.Sm_id,
		ord.Date_created,
		ord.Oof_shard,
		time.Now().Unix())

	if err != nil {
		return
	}

	query = `
		INSERT INTO delivery
		(order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err = tx.Exec(ctx,
		query,
		ord.Order_uid,
		ord.Delivery.Name,
		ord.Delivery.Phone,
		ord.Delivery.Zip,
		ord.Delivery.City,
		ord.Delivery.Address,
		ord.Delivery.Region,
		ord.Delivery.Email)

	if err != nil {
		return
	}

	query = `
		INSERT INTO payment
			(order_uid, transaction, request_id, currency, provider, amount,
			payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
	_, err = tx.Exec(ctx,
		query,
		ord.Order_uid,
		ord.Payment.Transaction,
		ord.Payment.Request_id,
		ord.Payment.Currency,
		ord.Payment.Provider,
		ord.Payment.Amount,
		ord.Payment.Payment_dt,
		ord.Payment.Bank,
		ord.Payment.Delivery_cost,
		ord.Payment.Goods_total,
		ord.Payment.Custom_fee)

	if err != nil {
		return
	}

	query = `
		INSERT INTO items
			(order_uid, chrt_id, track_number, price, rid, name, sale, size,
			total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`

	for i, item := range ord.Items {
		_, err = tx.Exec(ctx,
			query,
			ord.Order_uid,
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
			item.Status)

		if err != nil {
			return
		}

		i++
	}

	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}

	return
}

//Получение новых данных из бд
func (db *DB) GetOrder(ctx context.Context, order_uid string) (ord *models.Orders, err error) {

	ord = new(models.Orders)

	query := `
		SELECT order_uid, track_number, entry, locale, internal_signature,
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders o
		WHERE o.order_uid = $1
	`
	log.Printf("get orders")

	err = db.pool.QueryRow(ctx, query, order_uid).Scan(
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
	log.Printf("%v", err)

	if err != nil {
		log.Printf("from orders:%v", err)
		return
	}

	query = `
		SELECT name, phone, zip, city, address, region, email 
		FROM delivery d
		WHERE d.order_uid = $1
	`
	log.Printf("get delivery")
	err = db.pool.QueryRow(ctx, query, order_uid).Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	)

	if err != nil {
		log.Printf("from delivery:%v", err)
		return
	}

	query = `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
		FROM payment p
		WHERE p.order_uid = $1;
	`
	log.Printf("get payment")
	err = db.pool.QueryRow(ctx, query, order_uid).Scan(
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

	if err != nil {
		log.Printf("from payment:%v", err)
		return
	}

	query = `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status 
		FROM items i
		WHERE i.order_uid = $1;
	`

	rows, err := db.pool.Query(ctx, query, order_uid)
	defer rows.Close()

	for rows.Next() {
		item := new(models.Item)
		log.Printf("get items")
		err = rows.Scan(
			&item.Chrt_id,
			&item.Track_number,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.Total_price,
			&item.Nm_id,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			log.Printf("from items:%v", err)
			return
		}
		ord.Items = append(ord.Items, *item)
	}

	return ord, nil

}

func (db *DB) GetInitialCache(ctx context.Context) (*cache.Cache, error) {
	cache := cache.NewCache()

	now := time.Now().Unix()
	forCache := make([]string, 0, 10)

	query := `
	SELECT order_uid
	FROM orders o
	WHERE o.time_of_creation > $1;
	`

	rows, err := db.pool.Query(ctx, query, now-60*60*24)
	defer rows.Close()
	log.Printf("Зашли   GetInitialCache 1: %v", err)
	//нахождение всех заказов
	for rows.Next() {
		log.Printf("Зашли   GetInitialCache 2: %v", err)
		var order_uid *string
		err = rows.Scan(
			&order_uid,
		)
		log.Printf("Звшли 4: %v", err)
		if err != nil {
			log.Printf("failed cache %v", err)
			return cache, err
		}
		forCache = append(forCache, *order_uid)
		log.Printf("forCache: %v", forCache)
	}

	//get и marshal одной записи для кеша
	for _, order_uid := range forCache {

		ord, err := db.GetOrder(ctx, order_uid)

		message, err := json.Marshal(ord)
		if err != nil {
			log.Printf("Marshal for cache <%v>", err)
			return cache, err
		}
		cache.Data[order_uid] = string(message)
	}

	return cache, nil
}

//`SELECET orders max(time) - mid(time)`
/*
for rows.Next() {
	u := models.User{}
	err = rows.Scan(
		&u.Id,
		&u.Name,
		&u.Rating,
		&u.Rank,
		&u.Streak,
	)

	if err != nil {
		return
	}

	board = append(board, u)
}
*/
/*
	db.pool.Exec(`INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`)

	   	query := `
	   	INSERT INTO orders
	   		(order_uid, track_number, entry, locale, internal_signature,customer_id,
	   		delivery_service, shardkey, sm_id, date_created, oof_shard, time_of_creation)
	   	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	   `

	   	_, err = db.pool.Exec(ctx, query,
	   		ord.Order_uid,
	   		ord.Track_number,
	   		ord.Entry,
	   		ord.Locale,
	   		ord.Internal_signature,
	   		ord.Customer_id,
	   		ord.Delivery_service,
	   		ord.Shardkey,
	   		ord.Sm_id,
	   		ord.Date_created,
	   		ord.Oof_shard,
	   		time.Now().Unix())

	   	log.Printf("err exec %v", err)
*/

//http://localhost:8080/postgres?id=b563feb7b2b84b6test
