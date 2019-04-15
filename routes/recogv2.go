package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func RecogV2() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		fileName, err := saveFileFromReq(r)
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(fileName)

		encodingsPath, err := filepath.Abs("./encodings.pickle")
		if err != nil {
			log.Fatal(err)
		}
		unknownImageFile, err := filepath.Abs(fileName)
		cmd := exec.Command("python3", "-m", "recognize_faces_image.py", "-e", encodingsPath, "-i", unknownImageFile, "-d", "hog")
		out, _ := cmd.Output()
		if string(out) == "" {
			w.Write([]byte("no results found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		matches := []Match{}
		for _, x := range strings.Split(string(out), "\n") {
			if x == "" {
				continue
			}
			matches = append(matches, Match{Name: x})
		}
		if len(matches) == 0 {
			w.Write([]byte("no results found"))
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(struct{ PossibleMatches []Match }{PossibleMatches: matches})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logrus.Info("Responded.")
	}
}
