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
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

var msgsToIgnore = []string{
	"no_persons_found",
	"Ignoring file",
	"unknown_person",
	"More than one face found",
}

const excluder = ""
const unknownPerson = ""
const moreThanOneFace = ""

// Output Output

// Threshold Threshold
const Threshold = 5

// CurrDir CurrDir
func CurrDir() ([]string, error) {

	info, err := ioutil.ReadDir("./known")
	if err != nil {
		return nil, err
	}
	listDirs := []string{}
	for _, f := range info {
		listDirs = append(listDirs, f.Name())
	}
	return listDirs, nil
}

func ignore(s string) bool {
	for _, x := range msgsToIgnore {
		if strings.Contains(s, x) || s == "" {
			return true
		}

	}
	return false
}

// Recognize Recognize
func Recognize() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		savedFilePath, err := saveFileFromReq(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer os.Remove(savedFilePath)
		dirs, err := CurrDir()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		done := 0
		respCh := make(chan string)
		errCh := make(chan error)
		for i, d := range dirs {
			spew.Dump(i)
			go func(d string, i int) {
				matches := []string{}
				t := time.Now()
				defer func() {
					logrus.Infof("goroutine %d took %f seconds", i, time.Since(t).Seconds())
				}()
				cmd := exec.Command("python3", "-m", "face_recognition.face_recognition_cli", "./known/"+d, "./unknown")
				out, err := cmd.Output()
				if err != nil {
					errCh <- err
					return
				}

				for _, x := range strings.Split(string(out), "\n") {
					if ignore(x) {
						continue
					}
					matches = append(matches, strings.Split(x, ",")[1])
				}
				if len(matches) >= Threshold {
					respCh <- d
				}
				respCh <- ""
			}(d, i)
		}
		response := []string{}
		for {
			select {
			case name := <-respCh:
				done++
				spew.Dump(done)
				if done >= len(dirs) {
					w.Header().Set("Content-Type", "application/json")
					err := json.NewEncoder(w).Encode(struct{ PossibleMatches []string }{PossibleMatches: response})
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
					logrus.Info("Responded.")
					return
				}
				if name == "" {
					continue
				}
				alreadyExists := false
				for _, x := range response {
					if x == name {
						alreadyExists = true
					}
				}
				if !alreadyExists {
					response = append(response, name)
				}
			case err := <-errCh:
				w.Write([]byte(err.Error()))
				return
			}
		}

	}
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
