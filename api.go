package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func NewAPIServer(listenAddr string, store Storage) *APIserver {
	return &APIserver{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIserver) Run() {
	router := mux.NewRouter()
	
	router.HandleFunc("/accounts/{id}", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/accounts", makeHttpHandleFunc(s.handleAccounts))

	log.Println("JSON API Server Running on Port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIserver) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	
	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountById(w, r)
	}
	if r.Method == "DELETE" {
		_, err := s.handleDeleteAccountById(w, r)

		if err != nil {
			return err
		}
	}

	return fmt.Errorf("method not allowed: %s", r.Method)
}

func (s *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	strId := mux.Vars(r)["id"]
	id, err := convertId(strId)

	if err != nil {
		return err
	}
	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIserver) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err :=  s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	// &Account{}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	
	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIserver) handleDeleteAccountById(w http.ResponseWriter, r *http.Request) (bool, error) {
	strId := mux.Vars(r)["id"]
	id, err := convertId(strId)

	if err != nil {
		return false, err
	}

	res, err := s.store.DeleteAccount(id)

	if err != nil {
		return false, err
	}
	return res, nil
}

func (s *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIserver struct {
	listenAddr string
	store Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func convertId(strId string) (int, error) {
	id, err := strconv.Atoi(strId)

	if err != nil {
		return 0, err
	}

	return id, nil
}