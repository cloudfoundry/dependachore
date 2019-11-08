package dependachore

import (
	"dependachore/tracker"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TrackerActivity struct {
	Changes []Change `json:"changes"`
	Kind    string   `json:"kind"`
}

type Change struct {
	Kind      string        `json:"kind"`
	Type      string        `json:"change_type"`
	NewValues tracker.Story `json:"new_values"`
}

//go:generate counterfeiter . TrackerClient
type TrackerClient interface {
	Get(storyID int) (tracker.Story, error)
	MoveAndChorify(storyID, afterStoryID int) error
}

type Handler struct {
	trackerClient   TrackerClient
	releaseMarkerID int
}

func NewHandler(trackerClient TrackerClient, releaseMarkerID int) Handler {
	return Handler{trackerClient: trackerClient, releaseMarkerID: releaseMarkerID}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		writeError(w, http.StatusMethodNotAllowed, "Sorry, only POST methods are supported.")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read body: %v", err)
		return
	}
	defer r.Body.Close()

	activity := TrackerActivity{}
	err = json.Unmarshal(body, &activity)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to unmarshal body to tracker activity: %v\n%s\n", err, string(body))
		return
	}

	if storyID, isDependabotActivity := extractStoryIDFromDependabotActivity(activity); isDependabotActivity {
		_ = h.trackerClient.MoveAndChorify(storyID, h.releaseMarkerID)
	}
}

func writeError(w http.ResponseWriter, status int, msg string, subs ...interface{}) {
	http.Error(w, fmt.Sprintf(msg, subs...), status)
	fmt.Printf(msg, subs...)
}

func extractStoryIDFromDependabotActivity(activity TrackerActivity) (int, bool) {
	if activity.Kind != "story_create_activity" {
		return 0, false
	}

	for _, change := range activity.Changes {
		if change.Kind != "story" {
			continue
		}
		if strings.Contains(change.NewValues.Description, "@dependabot-preview[bot]") {
			return change.NewValues.ID, true
		}
	}

	return 0, false
}

