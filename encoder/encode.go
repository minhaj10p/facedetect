package encoder

import (
	"log"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Encode Encode
func Encode(datasetPath, encodingPath string) error {

	logrus.Info("Encoding face images")
	cmd := exec.Command("python3", "encode_faces.py", "-i", datasetPath, "-e", encodingPath, "-d", "hog")
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}
	return cmd.Wait()

}
