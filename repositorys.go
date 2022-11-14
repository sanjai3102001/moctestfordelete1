//
//  go-unit-test-sql
//

package main

// Repository represent the repositories
type Repository interface {
	Close()
	FindByID(id string) (*UserModel, error)
	Find() ([]*UserModel, error)
	Create(user *UserModel) error
	Update(user *UserModel) error
	Delete(id string) error
}

type UserModel struct {
	ID    string
	Name  string
	Email string
	Phone string
}
