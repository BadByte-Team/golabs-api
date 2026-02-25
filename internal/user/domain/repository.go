package userdomain

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)

	// SearchByUsername returns users whose username contains the query string (case-insensitive).
	SearchByUsername(query string) ([]*User, error)

	// List returns a paginated list of all users and the total count.
	List(offset, size int) ([]*User, int, error)

	Update(user *User) error
	UpdatePassword(userID string, passwordHash string) error

	UpdateRole(userID string, role string) error
	UpdatePoints(userID string, points int) error

	Ban(userID string) error
	Unban(userID string) error
}
