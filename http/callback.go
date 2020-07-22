package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Action struct {
	Type   int
	userid string
}

type History struct {
	serverVersion string
	changes       []Change
}

type Change struct {
}

type Response struct {
	key         string
	status      int
	url         string
	changesurl  string
	history     History
	users       []string
	actions     []Action
	lastsave    string
	notmodified bool
}

var callbackHandler = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {
	fmt.Println(r.Body)
	decoder := json.NewDecoder(r.Body)

	var res Response
	err := decoder.Decode(&res)

	fmt.Println(err)
	fmt.Println(res)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]int{"error": 0}

	json.NewEncoder(w).Encode(response)

	return 0, nil
})
