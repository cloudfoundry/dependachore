package dependachore

import (
	"dependachore/dependachore"
	"dependachore/tracker"
	"net/http"
	"os"
	"strconv"
)

func AskDependachore(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("API_KEY")
	project := os.Getenv("PROJECT_ID")
	projectID, err := strconv.Atoi(project)
	if err != nil {
		http.Error(w, "not a numeric project id", http.StatusBadRequest)
	}
	marker := os.Getenv("RELEASE_MARKER_ID")
	markerID, err := strconv.Atoi(marker)
	if err != nil {
		http.Error(w, "not a numeric release marker id", http.StatusBadRequest)
	}

	trackerClient := tracker.NewClient(apiKey, projectID)
	dependachore.NewHandler(trackerClient, markerID).Handle(w, r)
}
