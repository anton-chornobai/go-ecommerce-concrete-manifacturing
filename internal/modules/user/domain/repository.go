package domain

type Repository interface {
	Create(user UserCreated) error
	GetByPhone(number string) (User, error)
	GetByEmail(email string) (User, error)
	SignUpByEmail(user *UserCreatedWithEmail) (error)
}
