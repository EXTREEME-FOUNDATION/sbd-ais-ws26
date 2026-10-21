package model

import "time"

type Order struct {
	DrinkID    uint64    `json:"drink_id"` // foreign key
	Created_at time.Time `json:"created_at"`
	Amount     uint64    `json:"amount"`
	// todo Add fields: CreatedAt (time.Time), Amount with suitable types
	// todo json attributes need to be snakecase, i.e. name, created_at, my_variable, ..
}

/*
Order
{
    amount      integer
    created_at  string
    drink_id    integer // foreign key
}
*/
