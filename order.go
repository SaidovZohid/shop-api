package main

type OrderParam struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type Order struct {
	Id          int64        `json:"id"`
	CustomerID  int64        `json:"customer_id"`
	Address     string       `json:"address"`
	TotalAmount float64      `json:"total_amount"`
	OrderItems  []*OrderItem `json:"order_items"`
}

type OrderItem struct {
	Id          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ProductName string  `json:"product_name"`
	ProductID   int64   `json:"product_id"`
	Count       int64   `json:"count"`
	TotalPrice  float64 `json:"total_price"`
	Status      bool    `json:"status"`
}

type AllOrders struct {
	Orders []*Order `json:"orders"`
	Count int64 `json:"count"`
}

func (d *DBManager) CreateOrder(o *Order) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return 0, err
	}
	queryCheckBalance := `
		UPDATE customer SET balance = balance - $1 WHERE id = $2
	`
	_, err = tx.Exec(queryCheckBalance, o.TotalAmount, o.CustomerID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	queryAddOrder := `
		INSERT INTO orders(
			customer_id,
			address,
			total_amount
		) VALUES ($1, $2, $3) RETURNING id
	`
	queryAddItems := `
			INSERT INTO order_items(
				order_id,
				product_name,
				product_id,
				count,
				total_price,
				status
			) VALUES ($1, $2, $3, $4, $5, $6)
	`
	err = tx.QueryRow(queryAddOrder, o.CustomerID, o.Address, o.TotalAmount).Scan(&o.Id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, v := range o.OrderItems {
		_, err := tx.Exec(
			queryAddItems,
			o.Id,
			v.ProductName,
			v.ProductID,
			v.Count,
			v.TotalPrice,
			v.Status,
		)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return o.Id, nil
}

func (d *DBManager) GetOrder(order_id int64) (*Order, error) {
	queryGetOrder := `
		SELECT 
			id,
			customer_id,
			address,
			total_amount
		FROM orders WHERE id = $1
	`
	queryGetItems := `
		SELECT 
			id,
			order_id,
			product_name,
			product_id,
			count,
			total_price,
			status
		FROM order_items WHERE order_id = $1
	`
	var result Order
	err := d.db.QueryRow(
		queryGetOrder,
		order_id,
	).Scan(
		&result.Id,
		&result.CustomerID,
		&result.Address,
		&result.TotalAmount,
	)
	if err != nil {
		return nil, err
	}
	rows, err := d.db.Query(queryGetItems, order_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order_item OrderItem
		err := rows.Scan(
			&order_item.Id,
			&order_item.OrderID,
			&order_item.ProductName,
			&order_item.ProductID,
			&order_item.Count,
			&order_item.TotalPrice,
			&order_item.Status,
		)
		if err != nil {
			return nil, err
		}
		result.OrderItems = append(result.OrderItems, &order_item)
	}
	return &result, nil
}

func (d *DBManager) UpdateOrder(o *Order) (*Order, error) {
	var result Order
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	queryCheckBalance := `
		UPDATE customer SET balance = balance - $1 WHERE id = $2
	`
	_, err = tx.Exec(queryCheckBalance, o.TotalAmount, o.CustomerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	queryUpdateOrder := `
		UPDATE orders SET 
			customer_id = $1,
			address = $2,
			total_amount = $3
		WHERE id = $4 RETURNING id, customer_id, address, total_amount
	`
	err = tx.QueryRow(
		queryUpdateOrder, 
		o.CustomerID,
		o.Address,
		o.TotalAmount,
		o.Id,
	).Scan(
		&result.Id,
		&result.CustomerID,
		&result.Address,
		&result.TotalAmount,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	queryDeleteOrderItems := `
		DELETE FROM order_items WHERE order_id = $1
	`
	_, err = tx.Exec(queryDeleteOrderItems, o.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	queryAddItems := `
		INSERT INTO order_items(
			order_id,
			product_name,
			product_id,
			count,
			total_price,
			status
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, order_id, product_name, product_id, count, total_price, status
	`
	for _, v := range o.OrderItems {
		var order_item OrderItem
		err = tx.QueryRow(
			queryAddItems,
			result.Id,
			v.ProductName,
			v.ProductID,
			v.Count,
			v.TotalPrice,
			v.Status,
		).Scan(
			&order_item.Id,
			&order_item.OrderID,
			&order_item.ProductName,
			&order_item.ProductID,
			&order_item.Count,
			&order_item.TotalPrice,
			&order_item.Status,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		result.OrderItems = append(result.OrderItems, &order_item)
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &result, nil
}

func (d *DBManager) DeleteOrder(order_id int64) error {
	queryDeleteItem := `DELETE FROM order_items WHERE order_id = $1`
	queryDeleteOrder := `DELETE FROM orders WHERE id = $1`
	_, err := d.db.Exec(queryDeleteItem, order_id)
	if err != nil {
		return err
	}
	_, err = d.db.Exec(queryDeleteOrder, order_id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DBManager) GetAllOrders(params *OrderParam) (*AllOrders, error) {
	var allorders AllOrders
	filter := "	WHERE true "
	offset := (params.Page - 1) * params.Limit
	queryGetAll := `
		SELECT
			id,
			customer_id,
			address,
			total_amount
		FROM orders
	` + filter + " LIMIT $1 OFFSET $2"
	queryGetItems := `
		SELECT 
			id,
			order_id,
			product_name,
			product_id,
			count,
			total_price,
			status
		FROM order_items WHERE order_id = $1
	`
	rows, err := d.db.Query(queryGetAll, params.Limit, offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result Order
		err := rows.Scan(
			&result.Id,
			&result.CustomerID,
			&result.Address,
			&result.TotalAmount,
		)
		if err != nil {
			return nil, err
		}
		rows2, err := d.db.Query(queryGetItems, result.Id)
		for rows2.Next() {
			var orderItem OrderItem
			err := rows2.Scan(
				&orderItem.Id,
				&orderItem.OrderID,
				&orderItem.ProductName,
				&orderItem.ProductID,
				&orderItem.Count,
				&orderItem.TotalPrice,
				&orderItem.Status,
			)
			if err != nil {
				return nil, err
			}
			result.OrderItems = append(result.OrderItems, &orderItem)
		}
		allorders.Orders = append(allorders.Orders, &result)
	}
	queryCount := "SELECT count(1) FROM orders"
	err = d.db.QueryRow(queryCount).Scan(&allorders.Count)
	if err != nil {
		return nil, err
	}
	return &allorders, nil
}
