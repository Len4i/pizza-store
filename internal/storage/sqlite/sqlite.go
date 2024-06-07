package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

type Order struct {
	Size      string `json:"size"`
	Amount    int    `json:"amount"`
	PizzaType string `json:"pizza-type"`
}

func New(storagePath string) (*Storage, error) {

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		size TEXT NOT NULL,
		amount NUMBER NOT NULL,
		pizza_type TEXT NOT NULL);
	`)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveOrder(order Order) (int64, error) {

	stmt, err := s.db.Prepare("INSERT INTO orders(size, amount, pizza_type) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(order.Size, order.Amount, order.PizzaType)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetOrder(id int64) (Order, error) {

	stmt, err := s.db.Prepare("SELECT size, amount, pizza_type FROM orders WHERE id = ?")
	if err != nil {
		return Order{}, err
	}

	var order Order

	err = stmt.QueryRow(id).Scan(&order.Size, &order.Amount, &order.PizzaType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, fmt.Errorf("order with id %d not found", id)
		}

		return Order{}, err
	}

	return order, nil
}

// TODO: implement deleteOrder
