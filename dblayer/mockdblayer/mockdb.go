package mockdblayer

import (
	"errors"
	"log"
	"strconv"

	persistence "github.com/omgitsotis/user-service/dblayer/persistence"
)

type MockEventEmitter struct{}

func (mee *MockEventEmitter) emitEvent(eventType, msg string) {
	log.Printf("emiting %s event with msg %s\n", eventType, msg)
}

type MockDatabase struct {
	Users        []*persistence.User
	EventEmitter MockEventEmitter
	IDCount      int
}

func NewMockDatabase() *MockDatabase {
	users := make([]*persistence.User, 0)
	mockEmitter := MockEventEmitter{}
	return &MockDatabase{users, mockEmitter, 1}
}

func (db *MockDatabase) AddUser(user persistence.User) (*persistence.User, error) {
	user.ID = strconv.Itoa(db.IDCount)
	db.IDCount++
	db.Users = append(db.Users, &user)

	log.Printf("[MockDB] added new user %s\n", user.ID)
	db.EventEmitter.emitEvent("user created", "mock-user-json")

	return &user, nil
}

func (db *MockDatabase) FindUserByID(id string) (*persistence.User, error) {
	for _, user := range db.Users {
		if user.ID == id {
			log.Printf("[MockDB] found user %s\n", user.ID)
			return user, nil
		}
	}

	return nil, errors.New("no user found with ID")
}

func (db *MockDatabase) DeleteUser(id string) error {
	indexToDelete := -1
	for i, user := range db.Users {
		if user.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return errors.New("no user found with ID")
	}

	db.Users = append(db.Users[:indexToDelete], db.Users[indexToDelete+1:]...)
	db.EventEmitter.emitEvent("user deleted", "mock-user-json")
	return nil
}

func (db *MockDatabase) FindUserByCriteria(criteria string, value string) ([]*persistence.User, error) {
	results := make([]*persistence.User, 0)

	for _, user := range db.Users {
		switch criteria {
		case "country":
			if user.Country == value {
				results = append(results, user)
			}
		case "first_name":
			if user.FirstName == value {
				results = append(results, user)
			}
		case "last_name":
			if user.FirstName == value {
				results = append(results, user)
			}
		case "nickname":
			if user.FirstName == value {
				results = append(results, user)
			}
		case "email":
			if user.FirstName == value {
				results = append(results, user)
			}
		default:
			return nil, errors.New("invalid search criteria")
		}
	}

	return results, nil
}

func (db *MockDatabase) UpdateUser(u persistence.User) (*persistence.User, error) {
	for _, user := range db.Users {
		if user.ID == u.ID {
			if u.FirstName != "" {
				user.FirstName = u.FirstName
			}

			if u.LastName != "" {
				user.LastName = u.LastName
			}

			if u.Country != "" {
				user.Country = u.Country
			}

			if u.Nickname != "" {
				user.Nickname = u.Nickname
			}

			if u.Email != "" {
				user.Email = u.Email
			}

			if u.Password != "" {
				user.FirstName = u.Password
			}
		}

		return user, nil
	}

	db.EventEmitter.emitEvent("user deleted", "mock-user-json")
	return nil, errors.New("invalid search criteria")
}
