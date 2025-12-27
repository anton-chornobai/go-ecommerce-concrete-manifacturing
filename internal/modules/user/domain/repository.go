package domain

type Repository interface {
	Create(user User) (error)
	GetByPhone(number string) (User, error)
}	
