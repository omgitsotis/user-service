package client

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dblayer "github.com/omgitsotis/user-service/dblayer"
	"github.com/omgitsotis/user-service/dblayer/persistence"
)

// userServiceHandler is the handler for the routes of the user handler. It has
// one field, a database layer interface to handle calls to the database.
type userServiceHandler struct {
	dbHandler dblayer.DatabaseHandler
}

// newUserHandler creates a new userServiceHandler with a provided database
// lasyer
func newUserHandler(dbh dblayer.DatabaseHandler) *userServiceHandler {
	return &userServiceHandler{dbHandler: dbh}
}

// ServeAPI creates the http router and the routes for the user service
func ServeAPI(dbh dblayer.DatabaseHandler, endpoint string) error {
	client := newUserHandler(dbh)
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(client.healthcheck)

	r.Methods("GET").Path("/user/{id}").HandlerFunc(client.getUserHandler)
	r.Methods("PUT").Path("/user/{id}").HandlerFunc(client.updateUserHandler)
	r.Methods("POST").Path("/user").HandlerFunc(client.addUserHandler)
	r.Methods("DELETE").Path("/user/{id}").HandlerFunc(client.deleteUserHandler)

	// I was contenplating using query params for this route, but I was not sure
	// if multiple params were allowed, so I chose a more rigid option here.
	r.Methods("GET").Path("/search/{criteria}/{search}").HandlerFunc(client.searchUserHandler)
	log.Println("[UserServiceHandler] Server started")
	return http.ListenAndServe(endpoint, r)
}

// ErrorResponse is a json object to hold error messages.
type ErrorResponse struct {
	Error string `json:"error"`
}

// getUserHandler takes an id and returns a user from the database
func (ush *userServiceHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[UserServiceHandler] Recieved GET request on /user")

	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		ush.writeErrorResponse(w, "no id found")
		return
	}

	user, err := ush.dbHandler.FindUserByID(userID)
	if err != nil {
		log.Printf("[UserServiceHandler] Error getting user: %s\n", err.Error())
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(user)
}

// addUserHandler takes all the inputs from the post form and creates a new
// user in the database
func (ush *userServiceHandler) addUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[UserServiceHandler] Recieved POST request on /user")

	r.ParseForm()
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	nickname := r.FormValue("nickname")
	password := r.FormValue("password")
	email := r.FormValue("email")
	country := r.FormValue("country")

	// Validation of the inputs would be here. Not sure if any are mandatory
	// fields at this point

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
		log.Printf("[UserServiceHandler] Error adding new user: %s\n", err.Error())
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(addedUser)
}

// getUserHandler takes an id and removes it from the database
func (ush *userServiceHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[UserServiceHandler] Recieved DELETE request on /user")
	vars := mux.Vars(r)
	userID, ok := vars["id"]
	if !ok {
		log.Println("[UserServiceHandler] no id found in path")
		ush.writeErrorResponse(w, "no id found")
		return
	}

	if err := ush.dbHandler.DeleteUser(userID); err != nil {
		log.Printf("[UserServiceHandler] Error deleting user: %s\n", err.Error())
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

// searchUserHandler takes a criteria string and a search term string and returns
// a list of all the users that match that criteria
func (ush *userServiceHandler) searchUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[UserServiceHandler] Recieved GET request on /search")

	vars := mux.Vars(r)
	searchItem, ok := vars["search"]
	if !ok {
		log.Println("[UserServiceHandler] no search item found in path")
		ush.writeErrorResponse(w, "no search item found")
		return
	}

	criteria, ok := vars["criteria"]
	if !ok {
		log.Println("[UserServiceHandler] no criteria item found in path")
		ush.writeErrorResponse(w, "no criteria found")
		return
	}

	users, err := ush.dbHandler.FindUserByCriteria(criteria, searchItem)
	if err != nil {
		log.Printf("[UserServiceHandler] Error searching for users: %s\n", err.Error())
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
		log.Println("[UserServiceHandler] no id found in path")
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
		log.Printf("[UserServiceHandler] Error updating user: %s\n", err.Error())
		ush.writeErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(updUser)
}

func (ush *userServiceHandler) healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Println("[UserServiceHandler] Recieved GET request on /")
	w.Write([]byte(`{status : ok}`))
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
