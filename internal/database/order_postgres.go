package database

import (
	"context"
	"errors"
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
	log.Printf("get")

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

	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("%v", err)
	} else {
		return
	}

	query = `
		SELECT name, phone, zip, city, address, region, email 
		FROM delivery d
		WHERE o.order_uid = $1
	`
	log.Printf("get2")
	err = db.pool.QueryRow(ctx, query, order_uid).Scan(
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	} else {
		return
	}

	query = `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee 
		FROM payment p
		WHERE o.order_uid = $1;
	`
	log.Printf("get3")
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
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	} else {
		return
	}

	query = `
		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status 
		FROM items i
		WHERE o.order_uid = $1;
	`

	rows, err := db.pool.Query(ctx, query, order_uid)
	defer rows.Close()

	for rows.Next() {
		item := new(models.Item)
		log.Printf("get4")
		err = rows.Scan(
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
			err = nil
		} else {
			return
		}
		ord.Items = append(ord.Items, *item)
	}

	return

}

func (db *DB) GetInitialCache(ctx context.Context) (*cache.Cache, error) {
	cache := cache.NewCache()

	now := time.Now().Unix()
	forCache := make([]string, 0, 1000)

	query := `
	SELECT order_uid
	FROM orders o
	WHERE o.time_of_creation > $1
	`

	rows, err := db.pool.Query(ctx, query, now-60*60*24)
	log.Printf("Зашли %v", err)
	defer rows.Close()
	log.Printf("Звшли 2: %v", err)
	//нахождение всех заказов
	for rows.Next() {
		log.Printf("Звшли 3: %v", err)
		var order_uid string
		err = rows.Scan(
			order_uid,
		)
		log.Printf("Звшли 4: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		} else {
			return cache, err
		}
		forCache = append(forCache, order_uid)
		log.Printf("forCache: %v", forCache)
	}

	//get и marshal одной записи для кеша
	for i, order_uid := range forCache {

		ord, err := db.GetOrder(ctx, order_uid)

		message, err := json.Marshal(ord)
		if err != nil {
			log.Printf("Marshal for cache <%v>", err)
		}
		cache.Data[order_uid] = string(message)
		i++
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
