package client

import (
	"testing"
	"net/http"
    "net/http/httptest"
    "encoding/json"

	dblayer "github.com/omgitsotis/user-service/dblayer"
	persistence "github.com/omgitsotis/user-service/dblayer/persistence"
)

func AddTestUser(mockDB dblayer.DatabaseHandler) {
	mockUser := persistence.User {
		FirstName: "Klay",
		LastName:  "Thompson",
		Nickname:  "Splash Brother",
		Password:  "password",
		Email:     "klay_thompson@mail.com",
		Country:   "usa",
	}

	mockDB.AddUser(mockUser)
}

func TestGetUserHandler(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)
	req, err :=  http.NewRequest("GET", "/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Router(mockDB).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    
    var user persistence.User
    if err = json.NewDecoder(rr.Body).Decode(&user); err != nil {
    	t.Fatal(err)
    }

    if user.FirstName != "Klay" {
    	t.Errorf("handler returned wrong first name: got %v want %v",
    		user.FirstName, "Klay")
    }

    if user.Nickname != "Splash Brother" {
    	t.Errorf("handler returned wrong nickname: got %v want %v",
    		user.FirstName, "Splash Brother")
    }

}

func TestGetUserHandlerNoUser(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	req, err :=  http.NewRequest("GET", "/user/001", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Router(mockDB).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusNotFound)
    }
}
