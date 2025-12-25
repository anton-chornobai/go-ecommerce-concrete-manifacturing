package domain

type Repository interface {
	Create(user User) (error)
}	
