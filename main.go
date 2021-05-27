package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var bank *Bank

type reqAccount struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balance   float64 `json:"balance"`
}
type reqMove struct {
	SenderId    uint    `json:"sender_id"`
	RecipientId uint    `json:"recipient_id"`
	Amount      float64 `json:"amount"`
}

type respId struct {
	Id uint `json:"id"`
}
type respBalance struct {
	Balance float64 `json:"balance"`
}

func accountHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		var account reqAccount
		decoder := json.NewDecoder(r.Body)
		var err error
		var id uint
		err = decoder.Decode(&account)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			break
		}
		id, err = bank.CreateAccount(account.FirstName, account.LastName, account.Balance)
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			break
		}
		w.WriteHeader(http.StatusCreated)
		respId := new(respId)
		respId.Id = id
		_ = json.NewEncoder(w).Encode(respId)
	default:
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}

func accountBalanceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			break
		}
		balance, err := bank.GetBalanceById(uint(id))
		if err != nil {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			break
		}
		respBalance := new(respBalance)
		respBalance.Balance = balance
		_ = json.NewEncoder(w).Encode(respBalance)
	default:
		http.Error(w, "501 Not Implemented", http.StatusNotImplemented)
	}
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		var reqMove reqMove
		decoder := json.NewDecoder(r.Body)
		var err error
		err = decoder.Decode(&reqMove)
		if err != nil {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			break
		}
		if (reqMove.SenderId == 0) || (reqMove.RecipientId == 0) {
			http.Error(w, "400 Bad Request", http.StatusBadRequest)
			break
		}
		err = bank.MoveCash(reqMove.SenderId, reqMove.RecipientId, reqMove.Amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			break
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	bank = NewBank()

	router := mux.NewRouter()
	router.HandleFunc("/account", accountHandler)
	router.HandleFunc("/account/{id:[0-9]+}/balance", accountBalanceHandler)
	router.HandleFunc("/payment", paymentHandler)
	http.Handle("/", router)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
