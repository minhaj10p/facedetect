package routes

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/darahayes/go-boom"

	"github.com/minhaj10p/facedetect/encoder"
)

func createFile(path string, fromFile multipart.File) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, fromFile); err != nil {
		return err
	}

	if _, err := f.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func AddFace() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		inputFile, inputFileHeaders, err := r.FormFile("fileupload")
		if err != nil {
			boom.BadRequest(w, err.Error())
			return
		}
		path, err := getOrMakeDirName(r.FormValue("name"))
		if err != nil {
			boom.BadRequest(w, err.Error())
			return
		}
		err = createFile(fmt.Sprintf("%s/%s", path, inputFileHeaders.Filename), inputFile)
		if err != nil {
			boom.Internal(w)
			return
		}

		ep, err := filepath.Abs("encodings.pickle")
		if err != nil {
			boom.NotFound(w, err)
			return
		}
		datasetPath, err := filepath.Abs("known")
		if err != nil {
			boom.NotFound(w, err)
			return
		}
		if err := encoder.Encode(datasetPath, ep); err != nil {
			boom.Internal(w)
		}

		w.Write([]byte("Image successfully uploaded"))
	}
}

func getOrMakeDirName(faceName string) (string, error) {
	absPath, err := filepath.Abs("./known")
	if err != nil {
		return "", err
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		fname := strings.Split(file.Name(), ".")[0]
		if fname == faceName {
			p, err := filepath.Abs("./known/" + fname)
			if err != nil {
				return "", err
			}
			return p, nil
		}
	}

	p, err := filepath.Abs("./known/" + faceName)
	if err := os.Mkdir(p, 0777); err != nil {
		return "", err
	}
	return p, nil
}
