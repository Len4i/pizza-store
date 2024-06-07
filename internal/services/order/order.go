package order

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/Len4i/pizza-store/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Orderer interface {
	SaveOrder(order sqlite.Order) (int64, error)
	GetOrder(id int64) (sqlite.Order, error)
	// deleteOrder(id int64) error
}

type OrderService struct {
	storage Orderer
}

func New(storage Orderer) *OrderService {
	return &OrderService{
		storage: storage,
	}
}

func (o *OrderService) Create(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var order sqlite.Order
	err := dec.Decode(&order)
	if err != nil {
		if errors.Is(err, io.EOF) {
			http.Error(w, "empty request body", http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate values
	if order.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if order.Size != "family" && order.Size != "personal" {
		http.Error(w, "size must be 'family' or 'personal'", http.StatusBadRequest)
		return
	}
	if order.PizzaType != "margherita" && order.PizzaType != "pugliese" && order.PizzaType != "marinara" {
		http.Error(w, "pizza type must be 'margherita', 'pugliese' or 'marinara'", http.StatusBadRequest)
		return
	}

	id, err := o.storage.SaveOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, id)
}

func (o *OrderService) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}
	order, err := o.storage.GetOrder(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, order)
}
