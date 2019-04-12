package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Kagami/go-face"
	"github.com/davecgh/go-spew/spew"
)

const excluder = "Ignoring file"

type Output struct {
	Status string
}

func Recognize(rec *face.Recognizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		savedFilePath, err := saveFileFromReq(r)
		if err != nil {
			//			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer removeFile(savedFilePath)

		cmd := exec.Command("/usr/bin/python3", "./face_recognition_cli.py", "./known", "./unknown")
		out, err := cmd.Output()
		if err != nil {
			spew.Dump(err)
			log.Fatal(err.Error())
		}
		var outb, errb bytes.Buffer

		cmd.Stderr = &errb
		cmd.Stdout = &outb
		matches := []Output{}
		for _, x := range strings.Split(string(out), "\n") {
			if strings.Contains(x, excluder) || x == "" {
				continue
			}
			matches = append(matches, Output{Status: strings.Split(x, ",")[1]})
		}

		b, err := json.Marshal(matches)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(b)

	}

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
