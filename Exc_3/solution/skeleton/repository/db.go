package repository

import (
	"ordersystem/model"
	"time"
)

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

// todo
func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	// drinks := ...
	drinks := []model.Drink{
		{ID: 1, Name: "Pina Colada", Price: 1.5, Description: "Copakabana"},
		{ID: 2, Name: "RedBull Vodka", Price: 1.4, Description: "yes"},
		{ID: 3, Name: "Apple Juice", Price: 2.0, Description: "Freshly squeezed orange juice for babies only"},
	}

	orders := []model.Order{
		{DrinkID: 1, Created_at: time.Now(), Amount: 2},
		{DrinkID: 2, Created_at: time.Now(), Amount: 1},
		{DrinkID: 3, Created_at: time.Now(), Amount: 3},
	}

	// Init orders slice with some test data

	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

// todo
func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	// key = DrinkID, value = Amount of orders
	totalledOrders := make(map[uint64]uint64)
	for _, order := range db.orders {
		totalledOrders[order.DrinkID] += order.Amount
	}

	// totalledOrders map[uint64]uint64
	return totalledOrders
}

func (db *DatabaseHandler) AddOrder(order *model.Order) {
	// todo
	// add order to db.orders slice
	db.orders = append(db.orders, *order)
}
