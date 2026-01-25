package domain

type OrderRepository interface {
	Save(*Order) error
	Orders(limit int) ([]Order, error)
	// OrderWithUserInfo(userId string) (Order, error)
}
