package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minhaj10p/facerecog/routes"
)

const dataDir = "./models"
const port = 8080
const excluder = "Ignoring file"

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/recognition", routes.Recognize()).Methods("POST")
	r.HandleFunc("/face", routes.AddFace()).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
