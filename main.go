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
	r.HandleFunc("/v1/recognition", routes.Recognize()).Methods("POST")
	r.HandleFunc("/v2/recognition", routes.RecogV2()).Methods("POST")

	r.HandleFunc("/v1/face", routes.AddFace()).Methods("POST")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
