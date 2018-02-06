package dblayer

import (
	"errors"
	persistence "github.com/omgitsotis/user-service/dblayer/persistence"
	mockDB "github.com/omgitsotis/user-service/dblayer/mockdblayer"
)

type DBType string

type DatabaseHandler interface {
	AddUser(persistence.User) 		   (*persistence.User, error)
	FindUserByID(string) 			   (*persistence.User, error)
	DeleteUser(string) 				   (error)
	FindUserByCriteria(string, string) ([]*persistence.User, error)
	UpdateUser(persistence.User) 	   (*persistence.User, error)
}

const (
	MOCKDB DBType = "mockdb"
)

func NewPersistenceLayer(options DBType, connection string) (DatabaseHandler, error) {
	switch options {
	case MOCKDB:
		return mockDB.NewMockDatabase(), nil
	}

	return nil, errors.New("unsuported database type")
}
