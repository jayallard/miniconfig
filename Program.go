package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hi there")

	var router = mux.NewRouter()
	router.HandleFunc("/variables", GetVariables).Methods("GET")
	router.HandleFunc("/secrets", GetSecretIds).Methods("GET")
	router.HandleFunc("/secrets/{secretId}", GetSecret).Methods("GET")
	router.HandleFunc("/secrets/{secretId}/model", GetSecretModel).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetVariables(w http.ResponseWriter, _ *http.Request) {
	variables, _ := GetVariablesFromSecretsManager()
	json.NewEncoder(w).Encode(variables)
}

func GetSecretIds(w http.ResponseWriter, _ *http.Request) {
	var secretIds []SecretId
	secretIds = append(secretIds, SecretId{SecretId: "a", Description: "a"})
	json.NewEncoder(w).Encode(secretIds)
}

func GetSecret(w http.ResponseWriter, r *http.Request) {
	// TODO: long term goal: return the actual secret value
	// if the user has access to.
	var vars = mux.Vars(r)
	w.WriteHeader(http.StatusNotFound)
	fmt.Println("requested secret id: " + vars["secretId"])
}

func GetSecretModel(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Println("requested model for secret id: " + vars["secretId"])
	w.Write([]byte("okay!"))
}

type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SecretId struct {
	SecretId    string `json:"secretId"`
	Description string `json:"description"`
}
