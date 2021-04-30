package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	dTypes "codeGen.com/app/dataTypes"
	"codeGen.com/app/generators"
	"github.com/gorilla/mux"
)

func createEntity(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var entity dTypes.Entity
	json.Unmarshal(reqBody, &entity)
	generators.CreateModel(entity)
	generators.CreateInterface(entity)
}
func updateEntity(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	var entity dTypes.Entity
	json.Unmarshal(reqBody, &entity)
	generators.UpdateModel(entity)
	generators.UpdateInterface(entity)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/create", createEntity).Methods("POST")
	myRouter.HandleFunc("/update", updateEntity).Methods("POST")
	fmt.Println("Listening on port 10000")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	handleRequests()
}
