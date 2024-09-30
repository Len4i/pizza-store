package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

type Storage struct {
	DB *sql.DB
}

type Order struct {
	Size      string `json:"size"`
	Amount    int    `json:"amount"`
	PizzaType string `json:"pizza-type"`
}

func New(db *sql.DB) *Storage {
	return &Storage{
		DB: db,
	}
}

func (s *Storage) Close() error {
	return s.DB.Close()
}

func (s *Storage) SaveOrder(order Order) (int64, error) {

	stmt, err := s.DB.Prepare("INSERT INTO orders(size, amount, pizza_type) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
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

	stmt, err := s.DB.Prepare("SELECT size, amount, pizza_type FROM orders WHERE id = ?")
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
