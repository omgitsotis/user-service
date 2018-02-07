package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	dblayer "github.com/omgitsotis/user-service/dblayer"
	persistence "github.com/omgitsotis/user-service/dblayer/persistence"
)

func AddTestUser(mockDB dblayer.DatabaseHandler) {
	mockUser := persistence.User{
		FirstName: "Klay",
		LastName:  "Thompson",
		Nickname:  "Splash Brother",
		Password:  "password",
		Email:     "klay_thompson@mail.com",
		Country:   "usa",
	}

	mockDB.AddUser(mockUser)
}

func AddSearchUsers(mockDB dblayer.DatabaseHandler) {
	mockUser := persistence.User{
		FirstName: "Klay",
		LastName:  "Thompson",
		Nickname:  "Splash Brother",
		Password:  "password",
		Email:     "klay_thompson@mail.com",
		Country:   "usa",
	}

	mockUser2 := persistence.User{
		FirstName: "Serge",
		LastName:  "Ibaka",
		Nickname:  "Iblocka",
		Password:  "dorwssap",
		Email:     "serge_ibaka@mail.com",
		Country:   "cameroon",
	}

	mockUser3 := persistence.User{
		FirstName: "Steph",
		LastName:  "Curry",
		Nickname:  "Chef Curry",
		Password:  "pdarsosw",
		Email:     "steph_curry@mail.com",
		Country:   "usa",
	}

	mockDB.AddUser(mockUser)
	mockDB.AddUser(mockUser2)
	mockDB.AddUser(mockUser3)
}

func TestGetUserHandler(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)
	req, err := http.NewRequest("GET", "/user/1", nil)
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

	req, err := http.NewRequest("GET", "/user/001", nil)
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

func TestAddUserHandler(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)

	form := url.Values{}
	form.Add("first_name", "Otis")
	form.Add("last_name", "Simon")
	form.Add("nickname", "omgitsotis")
	form.Add("password", "p4ssw0rd")
	form.Add("email", "otis_simon@mail.com")
	form.Add("country", "UK")

	req, err := http.NewRequest("POST", "/user", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	if user.ID != "2" {
		t.Errorf("handler returned wrong ID: got %v want %v",
			user.ID, "2")
	}

	if user.FirstName != "Otis" {
		t.Errorf("handler returned wrong first name: got %v want %v",
			user.FirstName, "Otis")
	}

	if user.LastName != "Simon" {
		t.Errorf("handler returned wrong last name: got %v want %v",
			user.LastName, "Simon")
	}

	if user.Nickname != "omgitsotis" {
		t.Errorf("handler returned wrong nickname: got %v want %v",
			user.FirstName, "omgitsotis")
	}

	if user.Password != "p4ssw0rd" {
		t.Errorf("handler returned wrong password: got %v want %v",
			user.Password, "p4ssw0rd")
	}

	if user.Email != "otis_simon@mail.com" {
		t.Errorf("handler returned wrong email: got %v want %v",
			user.Email, "otis_simon@mail.com")
	}

	if user.Country != "UK" {
		t.Errorf("handler returned wrong email: got %v want %v",
			user.Country, "UK")
	}
}

func TestUpdateUser(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)

	form := url.Values{}
	form.Add("email", "otis_simon@mail.com")
	form.Add("country", "UK")

	req, err := http.NewRequest("PUT", "/user/1", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

	if user.ID != "1" {
		t.Errorf("handler returned wrong ID: got %v want %v",
			user.ID, "1")
	}

	if user.FirstName != "Klay" {
		t.Errorf("handler returned wrong first name: got %v want %v",
			user.FirstName, "Klay")
	}

	if user.LastName != "Thompson" {
		t.Errorf("handler returned wrong last name: got %v want %v",
			user.LastName, "Simon")
	}

	if user.Nickname != "Splash Brother" {
		t.Errorf("handler returned wrong nickname: got %v want %v",
			user.FirstName, "Splash Brother")
	}

	if user.Password != "password" {
		t.Errorf("handler returned wrong password: got %v want %v",
			user.Password, "password")
	}

	if user.Email != "otis_simon@mail.com" {
		t.Errorf("handler returned wrong email: got %v want %v",
			user.Email, "otis_simon@mail.com")
	}

	if user.Country != "UK" {
		t.Errorf("handler returned wrong email: got %v want %v",
			user.Country, "UK")
	}
}

func TestUpdateUserNoUser(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)

	form := url.Values{}
	form.Add("email", "otis_simon@mail.com")
	form.Add("country", "UK")

	req, err := http.NewRequest("PUT", "/user/2", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	Router(mockDB).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestDeleteUser(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)

	req, err := http.NewRequest("DELETE", "/user/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Router(mockDB).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteUserNoUser(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddTestUser(mockDB)

	req, err := http.NewRequest("DELETE", "/user/2", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	Router(mockDB).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestSearchUsers(t *testing.T) {
	mockDB, err := dblayer.NewPersistenceLayer(dblayer.MOCKDB, "")
	if err != nil {
		t.Fatal(err)
	}

	AddSearchUsers(mockDB)
	r := Router(mockDB)

	req, err := http.NewRequest("GET", "/search/country/usa", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	users := make([]persistence.User, 0)
	if err = json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Errorf("handler returned wrong number of users: got %v want %v",
			len(users), 2)
	}

	req, err = http.NewRequest("GET", "/search/nickname/Iblocka", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err = json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatal(err)
	}

	if len(users) != 1 {
		t.Errorf("handler returned wrong number of users: got %v want %v",
			len(users), 1)
	}

	req, err = http.NewRequest("GET", "/search/country/uk", nil)
	if err != nil {
		t.Fatal(err)
	}

	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if err = json.NewDecoder(rr.Body).Decode(&users); err != nil {
		t.Fatal(err)
	}

	if len(users) != 0 {
		t.Errorf("handler returned wrong number of users: got %v want %v",
			len(users), 0)
	}

}
