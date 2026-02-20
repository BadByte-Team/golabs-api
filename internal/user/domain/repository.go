package userdomain

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)

	Update(user *User) error
	UpdatePassword(userID string, passwordHash string) error

	UpdateRole(userID string, role string) error
	UpdatePoints(userID string, points int) error

	Ban(userID string) error
	Unban(userID string) error
}
