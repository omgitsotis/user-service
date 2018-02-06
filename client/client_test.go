package client

import (
	"encoding/json"
	"testing"

	dblayer "github.com/omgitsotis/user-service/dblayer"
)



func TestGetUserHandler(t testing.T) {
	handler := newUserHandler(dblayer.NewPersistenceLayer(dblayer.MOCKDB, ""))
}
