package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/minhaj10p/facedetect/encoder"
	"github.com/minhaj10p/facedetect/routes"
)

const dataDir = "./models"
const port = 8080
const excluder = "Ignoring file"

func main() {

	ep, err := filepath.Abs("encodings.pickle")
	if err != nil {
		log.Fatal(err)
	}
	datasetPath, err := filepath.Abs("known")
	if err != nil {
		log.Fatal(err)
	}

	//if encodings don't exist.
	if _, err := os.Stat(ep); err != nil {
		t := time.Now()
		if err := encoder.Encode(datasetPath, ep); err != nil {
			log.Fatalf("could not encode dataset. err: %s", err)
		}
		logrus.Infof("encoding takes %f seconds", time.Since(t).Seconds())
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/recognition", routes.Recognize()).Methods("POST")
	r.HandleFunc("/v2/recognition", routes.RecogV2()).Methods("POST")

	r.HandleFunc("/v1/face", routes.AddFace()).Methods("POST")
	logrus.Infof("listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}
