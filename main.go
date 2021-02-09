package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Login struct contains user login data
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type loginHandlers struct {
	logins map[string]Login
}

func newLoginHandlers() *loginHandlers {
	return &loginHandlers{
		logins: map[string]Login{},
	}
}

func (h *loginHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyIn, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("expected content-type 'application/json',	 got '%s'", contentType)))
		return
	}

	var login Login
	err = json.Unmarshal(bodyIn, &login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	//add/update login in logins map
	h.logins[login.Username] = login

	// bodyOut, err := json.Marshal(login)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(err.Error()))
	// }

	// w.Header().Add("content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(bodyOut)
}

func main() {
	loginHandlers := newLoginHandlers()
	http.HandleFunc("/login", loginHandlers.post)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
