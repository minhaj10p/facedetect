package main

import (
	"fmt"
	"log"
	"net/http"

	face "github.com/Kagami/go-face"
	"github.com/gorilla/mux"
	"github.com/minhajuddinkhan/facerecog/routes"
)

const dataDir = "./models"
const port = 8080
const excluder = "Ignoring file"

func main() {

	// cmd := exec.Command("/usr/bin/python3", "./face_recognition_cli.py", "./known", "./unknown")
	// out, err := cmd.Output()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// var outb, errb bytes.Buffer

	// cmd.Stderr = &errb
	// cmd.Stdout = &outb
	// matches := []string{}
	// for _, x := range strings.Split(string(out), "\n") {
	// 	if strings.Contains(x, excluder) {
	// 		continue
	// 	}
	// 	matches = append(matches, x)
	// }
	// spew.Dump(matches)

	//	spew.Dump(errb.String(), outb.String())
	r := mux.NewRouter()
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Fatal(err)
	}
	r.HandleFunc("/recognition", routes.Recognize(rec))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r))
}

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
// func main2() {
// 	// Init the recognizer.
// 	// rec, err := face.NewRecognizer("./models")
// 	// if err != nil {
// 	// 	log.Fatalf("Can't init face recognizer: %v", err)
// 	// }
// 	// Free the resources when you're finished.
// 	//	defer rec.Close()

// 	// Test image with 10 faces.
// 	testImagePristin := filepath.Join(dataDir, "pristin.jpg")
// 	// Recognize faces on that image.
// 	faces, err := rec.RecognizeFile(testImagePristin)
// 	if err != nil {
// 		log.Fatalf("Can't recognize: %v", err)
// 	}
// 	if len(faces) != 10 {
// 		log.Fatalf("Wrong number of faces")
// 	}

// 	// Fill known samples. In the real world you would use a lot of images
// 	// for each person to get better classification results but in our
// 	// example we just get them from one big image.
// 	var samples []face.Descriptor
// 	var cats []int32
// 	for i, f := range faces {
// 		samples = append(samples, f.Descriptor)
// 		// Each face is unique on that image so goes to its own category.
// 		cats = append(cats, int32(i))
// 	}
// 	// Name the categories, i.e. people on the image.
// 	labels := []string{
// 		"Sungyeon", "Yehana", "Roa", "Eunwoo", "Xiyeon",
// 		"Kyulkyung", "Nayoung", "Rena", "Kyla", "Yuha",
// 	}
// 	// Pass samples to the recognizer.
// 	rec.SetSamples(samples, cats)

// 	// Now let's try to classify some not yet known image.
// 	testImageNayoung := filepath.Join(dataDir, "nayoung.jpg")
// 	nayoungFace, err := rec.RecognizeSingleFile(testImageNayoung)
// 	if err != nil {
// 		log.Fatalf("Can't recognize: %v", err)
// 	}
// 	if nayoungFace == nil {
// 		log.Fatalf("Not a single face on the image")
// 	}
// 	catID := rec.Classify(nayoungFace.Descriptor)
// 	if catID < 0 {
// 		log.Fatalf("Can't classify")
// 	}
// 	// Finally print the classified label. It should be "Nayoung".
// 	fmt.Println(labels[catID])
// }
