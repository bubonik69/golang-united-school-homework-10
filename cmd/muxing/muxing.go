package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	_"github.com/stretchr/testify/assert"
	"github.com/gorilla/mux"
)

//METHOD	REQUEST	RESPONSE
//+GET	/name/{PARAM}	body: Hello, PARAM!
//+GET	/bad	Status: 500
//POST	/data + Body PARAM	body: I got message:\nPARAM
//POST	/headers+ Headers{"a":"2", "b":"3"}	Header "a+b": "5"
func Start(host string, port int) {
	router := mux.NewRouter()
	router.HandleFunc("/name/{param}", paramPage).Methods("Get")
	router.HandleFunc("/bad", badPage).Methods("Get")
	router.HandleFunc("/data", dataPage).Methods("Post")
	router.HandleFunc("/headers", sumPage).Methods("Post")
	router.HandleFunc("/", notDefinedPage).Methods("Get")
	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), router); err != nil {
		log.Fatal(err)
	}

}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}
func notDefinedPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte("Hello"))
}

func paramPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	response := "Hello, " + vars["param"] + "!"
	w.Write([]byte(response))
}

func badPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func dataPage(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err == nil {
		response := "I got message:\n" + string(b)
		w.Write([]byte(response))
	}
}

func sumPage(w http.ResponseWriter, r *http.Request) {
	h := r.Header
	a := h["A"]
	b := h["B"]
	if a != nil && b != nil {
		a_int, err := strconv.Atoi(a[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		b_int, err := strconv.Atoi(b[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		sum := strconv.Itoa(a_int + b_int)
		w.Header().Set("a+b", sum)
		w.WriteHeader(http.StatusOK)
	}
}
