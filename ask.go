package dependachore

import (
	"encoding/base64"
	"fmt"
	"github.com/masters-of-cats/dependachore/dependachore"
	"github.com/masters-of-cats/dependachore/kms"
	"github.com/masters-of-cats/dependachore/tracker"
	"net/http"

	"os"
	"strconv"
)

var apiKey []byte

func AskDependachore(w http.ResponseWriter, r *http.Request) {
	if len(apiKey) == 0 {
		kmsProject := os.Getenv("ENC_PROJECT")
		location := os.Getenv("ENC_LOCATION")
		keyRing := os.Getenv("ENC_KEYRING")
		key := os.Getenv("ENC_KEY")

		kmsClient := kms.NewClient(kmsProject, location, keyRing, key)

		encKey := os.Getenv("ENC_API_KEY")
		encKeyBytes, err := base64.StdEncoding.DecodeString(encKey)
		if err != nil {
			fmt.Printf("can't decode base64 ENC_API_KEY: %v\n", err)
			http.Error(w, fmt.Sprintf("can't decode base64 ENC_API_KEY: %v", err), http.StatusInternalServerError)
			return
		}

		apiKey, err = kmsClient.Decrypt(encKeyBytes)
		if err != nil {
			fmt.Printf("can't decrypt api key: %v\n", err)
			http.Error(w, "can't decrypt api key", http.StatusInternalServerError)
			return
		}
	}

	project := os.Getenv("PROJECT_ID")
	projectID, err := strconv.Atoi(project)
	if err != nil {
		fmt.Println("not a numeric project id")
		http.Error(w, "not a numeric project id", http.StatusBadRequest)
		return
	}

	trackerClient := tracker.NewClient(string(apiKey), projectID)

	marker := os.Getenv("RELEASE_MARKER_ID")
	markerID, err := strconv.Atoi(marker)
	if err != nil {
		fmt.Println("not a numeric release marker id")
		http.Error(w, "not a numeric release marker id", http.StatusBadRequest)
		return
	}

	dependachore.NewHandler(trackerClient, markerID).Handle(w, r)
}
