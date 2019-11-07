package dependachore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type TrackerActivity struct {
	PrimaryResources []PrimaryResource `json:"primary_resources"`
	Changes          []Change          `json:"changes"`
	Kind             string            `json:"kind"`
	PerformedBy      Performer         `json:"performed_by"`
}

type PrimaryResource struct {
	StoryType string `json:"story_type"`
	Name      string `json:"name"`
	ID        int    `json:"id"`
}

type Change struct {
	NewValues Story `json:"new_values"`
}

type Story struct {
	Description   string `json:"description"`
	RequestedByID int    `json:"requested_by_id"`
	StoryType     string `json:"story_type"`
	BeforeID      int    `json:"before_id"`
	AfterID       int    `json:"after_id"`
	ID            int    `json:"id"`
}

type Performer struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func AskDependachore(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Sorry, only POST methods are supported.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, "Failed to read body: %v", err)
		return
	}
	defer r.Body.Close()

	activity := TrackerActivity{}
	err = json.Unmarshal(body, &activity)
	if err != nil {
		writeError(w, "Failed to unmarshal body to tracker activity: %v\n%s\n", err, string(body))
		return
	}

	fmt.Printf("activity = %+v\n", activity)
	fmt.Printf("is dependabot thing: %t\n", isDependabotThing(activity))
}

func writeError(w http.ResponseWriter, msg string, subs ...interface{}) {
	http.Error(w, fmt.Sprintf(msg, subs...), http.StatusInternalServerError)
	fmt.Printf(msg, subs...)
}

func isDependabotThing(activity TrackerActivity) bool {
	if len(activity.Changes) == 0 {
		return false
	}

	//TODO: check this is a creation

	change := activity.Changes[0]
	return strings.Contains(change.NewValues.Description, "@dependabot-preview[bot]")
}
