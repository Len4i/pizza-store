package order

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
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
	log     *slog.Logger
}

func New(storage Orderer, log *slog.Logger) *OrderService {
	return &OrderService{
		storage: storage,
		log:     log,
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
			o.log.Debug("empty request body")
			return
		}
		o.log.Debug("failed to decode request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate values
	if order.Amount <= 0 {
		o.log.Debug("amount must be greater than 0", "request", order)
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if order.Size != "family" && order.Size != "personal" {
		o.log.Debug("size must be 'family' or 'personal'", "request", order)
		http.Error(w, "size must be 'family' or 'personal'", http.StatusBadRequest)
		return
	}
	if order.PizzaType != "margherita" && order.PizzaType != "pugliese" && order.PizzaType != "marinara" {
		o.log.Debug("pizza type must be 'margherita', 'pugliese' or 'marinara'", "request", order)
		http.Error(w, "pizza type must be 'margherita', 'pugliese' or 'marinara'", http.StatusBadRequest)
		return
	}

	id, err := o.storage.SaveOrder(order)
	if err != nil {
		o.log.Error("failed to save order", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]int64{"order_id": id})
}

func (o *OrderService) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		o.log.Debug("invalid id parameter", "request", id, "error", err)
		http.Error(w, "invalid id parameter", http.StatusBadRequest)
		return
	}
	order, err := o.storage.GetOrder(id)
	if err != nil {
		o.log.Error("failed to get order", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, order)
}
