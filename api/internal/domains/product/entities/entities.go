package entities

type Product struct {
	ID       int64  `json:"id" db:"product_id" example:"1"`
	Name     string `json:"name" db:"name" example:"Эклеры \"Домашние\" (с домашним кремом) (ЛОРИ)"`
	Weight   int    `json:"weight" db:"weight" example:"1000"`
	Lifetime int64  `json:"lifetime" db:"lifetime" example:"86400"`
	Img      string `json:"img" db:"img" example:"https://example.com/img.jpg"`
}

type ProductInOrder struct {
	ID     int64 `json:"id" db:"product_id" example:"1"`
	Amount int   `json:"amount" db:"amount" example:"2"`
}

type Order struct {
	ID          int64            `json:"id" db:"order_id" example:"1"`
	Products    []ProductInOrder `json:"products" db:"products"`
	ProductList []Product        `json:"-" db:"-" example:"-"`
	Customer    Customer
	Manager     Manager
	Date        int64 `json:"date" db:"date" example:"1630000000"`
}

type Customer struct {
	Phone   string `json:"phone" db:"customer_phone" example:"+79991234567"`
	Name    string `json:"name" db:"customer_name" example:"Иван Иванов"`
	Address string `json:"address" db:"customer_address" example:"г. Москва, ул. Ленина, д. 1"`
}

type Manager struct {
	ID    int64  `json:"id" db:"manager_id" example:"1"`
	Phone string `json:"phone" db:"phone" example:"+79991234567"`
	Email string `json:"email" db:"email" example:"mail@gmail.com"`
	State int    `json:"state" db:"state" example:"1"`
}
