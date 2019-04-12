package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const excluder = "Ignoring file"

// Output Output

type JSON struct {
	People []People
}
type People struct {
	Name   string   `json:"name"`
	Photos []string `json:"photos"`
}

// Recognize Recognize
func Recognize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		savedFilePath, err := saveFileFromReq(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer removeFile(savedFilePath)

		cmd := exec.Command("python3", "-m", "face_recognition.face_recognition_cli", "./known", "./unknown")
		out, err := cmd.Output()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		matches := []string{}
		for _, x := range strings.Split(string(out), "\n") {
			if strings.Contains(x, excluder) || x == "" {
				continue
			}
			matches = append(matches, strings.Split(x, ",")[1])
		}

		x, _ := filepath.Abs("routes/reflect.json")
		b, err := ioutil.ReadFile(x)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		var ppl JSON
		if err := json.Unmarshal(b, &ppl); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		foundPPL := []string{}
		actor := FindActor(ppl.People, matches)

		found := false
		for _, x := range foundPPL {
			if x == actor {
				found = true
			}
		}
		if !found {
			foundPPL = append(foundPPL, actor)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(foundPPL); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

}

func FindActor(ppl []People, matchedPhotos []string) string {

	for _, p := range ppl {
		for _, x := range p.Photos {
			for _, photo := range matchedPhotos {
				if photo == x {

					return p.Name
				}
			}

		}
	}
	return ""

}

func removeFile(s string) error {
	return os.Remove(s)
}

func saveFileFromReq(r *http.Request) (string, error) {

	var buf bytes.Buffer
	file, header, err := r.FormFile("fileupload")
	if err != nil {
		return "", fmt.Errorf("cannot get file from request. err: %v", err)
	}
	defer file.Close()
	n := strings.Split(header.Filename, ".")

	if _, err := io.Copy(&buf, file); err != nil {
		return "", fmt.Errorf("cannot copy content to file. err: %v", err)
	}
	fileName := n[0] + time.Now().Format("20060102150405") + "." + n[1]
	filePath := "./unknown/" + fileName

	f, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("cannot create file. err: %v", err)
	}
	if _, err := f.Write(buf.Bytes()); err != nil {
		return "", fmt.Errorf("cannot write to file. err: %v", err)
	}
	return filePath, nil

}

func AddIfNotPresent(arr []string, s string) []string {
	for _, x := range arr {
		if x == s {
			continue
		}
		arr = append(arr, s)
	}
	return arr
}
