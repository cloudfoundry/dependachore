package dependachore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/masters-of-cats/dependachore/tracker"
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
		log("chorifying story %s", storyID)
		err = h.trackerClient.MoveAndChorify(storyID, h.releaseMarkerID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to move/chorify story: %v\n", err)
		}
	}
}

func writeError(w http.ResponseWriter, status int, msg string, subs ...interface{}) {
	http.Error(w, fmt.Sprintf(msg, subs...), status)
	fmt.Printf(msg, subs...)
}

func log(msg string, subs ...interface{}) {
	fmt.Printf(msg+"\n", subs...)
}

func extractStoryIDFromDependabotActivity(activity TrackerActivity) (int, bool) {
	if activity.Kind != "story_create_activity" {
		log("ignoring activity with kind %s", activity.Kind)
		return 0, false
	}

	for _, change := range activity.Changes {
		if change.Kind != "story" {
			log("ignoring change with kind %s", change.Kind)
			continue
		}
		if strings.Contains(change.NewValues.Description, "@dependabot-preview[bot]") {
			return change.NewValues.ID, true
		}
		log("description of the change does not match the dependabot pattern: %s", change.NewValues.Description)
	}

	return 0, false
}
