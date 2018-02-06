package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dblayer "github.com/omgitsotis/user-service/dblayer"
	"github.com/omgitsotis/user-service/dblayer/persistence"
)

type userServiceHandler struct {
	dbHandler dblayer.DatabaseHandler
}

func newUserHandler(dbh dblayer.DatabaseHandler) *userServiceHandler {
	return &userServiceHandler{dbHandler: dbh}
}

func ServeAPI(dbh dblayer.DatabaseHandler, endpoint string) error {
	client := newUserHandler(dbh)
	r := mux.NewRouter()
	r.Methods("GET").Path("/user/{id}").HandlerFunc(client.getUserHandler)
	r.Methods("PUT").Path("/user/{id}").HandlerFunc(client.updateUserHandler)
	r.Methods("POST").Path("/user").HandlerFunc(client.addUserHandler)
	r.Methods("DELETE").Path("/user/{id}").HandlerFunc(client.deleteUserHandler)

	r.Methods("GET").Path("/search/{criteria}/{search}").HandlerFunc(client.searchUserHandler)

	return http.ListenAndServe(endpoint, r)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ush *userServiceHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got request")
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		ush.writeErrorResponse(w, "no id found")
		return
	}

	user, err := ush.dbHandler.FindUserByID(userID)
	if err != nil {
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(user)
}

func (ush *userServiceHandler) addUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	nickname := r.FormValue("nickname")
	password := r.FormValue("password")
	email := r.FormValue("email")
	country := r.FormValue("country")

	user := persistence.User{
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		Password:  password,
		Email:     email,
		Country:   country,
	}

	addedUser, err := ush.dbHandler.AddUser(user)
	if err != nil {
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(addedUser)
}

func (ush *userServiceHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		ush.writeErrorResponse(w, "no id found")
		return
	}

	if err := ush.dbHandler.DeleteUser(userID); err != nil {
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ush *userServiceHandler) searchUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	searchItem, ok := vars["search"]
	if !ok {
		ush.writeErrorResponse(w, "no search item found")
		return
	}

	criteria, ok := vars["criteria"]
	if !ok {
		ush.writeErrorResponse(w, "no criteria found")
		return
	}

	users, err := ush.dbHandler.FindUserByCriteria(criteria, searchItem)
	if err != nil {
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(&users)
}

func (ush *userServiceHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved PUT request on route /user")
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		ush.writeErrorResponse(w, "no id found")
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	nickname := r.FormValue("nickname")
	password := r.FormValue("password")
	email := r.FormValue("email")
	country := r.FormValue("country")

	user := persistence.User{
		ID:        userID,
		FirstName: firstName,
		LastName:  lastName,
		Nickname:  nickname,
		Password:  password,
		Email:     email,
		Country:   country,
	}

	updUser, err := ush.dbHandler.UpdateUser(user)
	if err != nil {
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(updUser)
}

func (ush *userServiceHandler) writeErrorResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	er := ErrorResponse{msg}
	output, oErr := json.Marshal(er)
	if oErr != nil {
		http.Error(w, oErr.Error(), http.StatusBadRequest)
		return
	}

	w.Write(output)
}
