package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Login struct contains user login data
type Login struct {
	Hash  string `json:"Hash"`
	Token string `json:"token"`
}

//Response struct contains response data
type Response struct {
	Success     bool
	Message     string
	RedirectURL string
}

//LoginHandlers contains map for logins
type LoginHandlers struct {
	logins map[string]Login
}

func newLoginHandlers() *LoginHandlers {
	return &LoginHandlers{
		logins: map[string]Login{},
	}
}

//Validate validates user credentials
func validateCredentials(login Login) bool {

	if login.Hash == "EOCtK6aNq4iF67IjxyS3LIB3ymQb0/iP+T/ptOQaQX8=" {
		return true
	}

	return false
}

func validateRequest(w http.ResponseWriter, r *http.Request) bool {
	method := r.Method
	if method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return false
	}

	if method != "POST" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf("'%s' method not allowed.", method)))
		return false
	}

	contentType := r.Header.Get("content-type")
	if contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("expected content-type 'application/json',	 got '%s'", contentType)))
		return false
	}

	return true
}

func createResponse(valid bool) Response {
	var response Response

	if valid {
		response.Success = true
		response.Message = "Success"
		response.RedirectURL = "http://onecause.com"
	} else {
		response.Success = false
		response.Message = "bad username/password"
		response.RedirectURL = ""
	}

	return response
}

func (h *LoginHandlers) post(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type")

	if !validateRequest(w, r) {
		return
	}

	//get body from request
	bodyIn, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	//create login object
	var login Login
	err = json.Unmarshal(bodyIn, &login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	//validate credentials and create response using that result
	responseObject := createResponse(validateCredentials(login))

	//serialize response
	response, err := json.Marshal(responseObject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	//add/update login in logins map
	if responseObject.Success {
		h.logins[login.Hash] = login
	}

	//write response
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	//used for testing

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
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
