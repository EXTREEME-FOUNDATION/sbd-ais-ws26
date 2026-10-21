package rest

import (
	"encoding/json"
	"net/http"
	"ordersystem/model"
	"ordersystem/repository"

	"errors"
	"strconv"
	"time"

	"github.com/go-chi/render"
)

func returnError(w http.ResponseWriter, r *http.Request, err error, msg ...string) {
	render.Status(r, http.StatusBadRequest)
	if err == nil {
		err = errors.New("")
	}
	if len(msg) == 0 {
		render.JSON(w, r, "Super BAD error: "+err.Error())
	} else {
		render.JSON(w, r, msg[0]+": "+err.Error())
	}
}

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := db.GetDrinks()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, sl)
	}
}

// GetOrders		godoc
// @tags 			Order
// @Description 	Returns all orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Router 			/api/order/all [get]
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := db.GetOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, sl)
	}
}

// GetOrdersTotal	godoc
// @tags 			Order
// @Description 	Returns the total of all orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Router 			/api/order/totalled [get]
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := db.GetTotalledOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, sl)
	}
}

type ProtoOrder struct {
	DrinkID    uint64 `json:"drink_id"`
	Created_at string `json:"created_at"` //string to allow empty for now
	Amount     uint64 `json:"amount"`
}

// PostOrder 		godoc
// @tags 			Order
// @Description 	Adds an order to the db
// @Accept 			json
// @Param 			b body model.ProtoOrder true "Order"
// @Produce  		json
// @Success 		200
// @Failure     	400
// @Router 			/api/order [post]
func PostOrder(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// todo
		// declare empty order struct
		newOrder := &model.Order{}

		var proto ProtoOrder // The protoOrder struct is used to delay the parsing of the date until emptyness is checked.
		err := json.NewDecoder(r.Body).Decode(&proto)

		if err != nil {
			returnError(w, r, err)
			return
		}

		if proto.Created_at == "" || proto.Created_at == "string" {
			newOrder.Created_at = time.Now()
		} else {
			created_at, err := time.Parse(time.DateTime, proto.Created_at) // Datetime is "2006-01-02 15:04:05". (The month is 02 and the day is 01, which sux. #Murica)
			if err != nil {
				returnError(w, r, err, "Couldn't parse date.")
				return
			}
			newOrder.Created_at = created_at
		}
		newOrder.Amount = proto.Amount
		newOrder.DrinkID = proto.DrinkID

		if newOrder.Amount == 0 {
			returnError(w, r, err, "Amount can't be 0 you absolute fool")
			return
		}
		found := false
		for _, drink := range db.GetDrinks() {
			if drink.ID == newOrder.DrinkID {
				found = true
				break
			}
		}
		if !found {
			returnError(w, r, err, "Drink with ID "+strconv.FormatUint(newOrder.DrinkID, 10)+" not found")
			return
		}
		//newOrder.Created_at = time.Now()
		// handle error and render Status 400

		// add to db
		db.AddOrder(newOrder)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
	}
}
