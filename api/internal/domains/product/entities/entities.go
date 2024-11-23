package entities

type Product struct {
	ID       int64  `json:"id" db:"product_id"`
	Name     string `json:"name" db:"name"`
	Weight   int    `json:"weight" db:"weight"`
	Lifetime int64  `json:"lifetime" db:"lifetime"`
	Img      string `json:"img" db:"img"`
}

type ProductInOrder struct {
	ID     int64 `json:"id" db:"product_id"`
	Amount int   `json:"amount" db:"amount"`
}

type Order struct {
	ID       int64            `json:"id"`
	Products []ProductInOrder `json:"products"`
	Customer Customer
	Manager  Manager
}

type Customer struct {
	Phone   string `json:"phone" db:"customer_phone"`
	Name    string `json:"name" db:"customer_name"`
	Address string `json:"address" db:"customer_address"`
}

type Manager struct {
	Phone string `json:"phone" db:"phone"`
	Email string `json:"email" db:"email"`
}
