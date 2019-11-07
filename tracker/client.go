package tracker

import (
	"bytes"
	"dependachore"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Client struct {
	apiKey    string
	projectID int
	storyPath string
}

const trackerURL = "https://www.pivotaltracker.com"

func NewClient(apiKey string, projectID int) *Client {
	return &Client{
		apiKey:    apiKey,
		projectID: projectID,
		storyPath: fmt.Sprintf("services/v5/projects/%d/stories", projectID),
	}
}

func (c *Client) Get(storyID int) (dependachore.Story, error) {
	request, err := c.createRequest(http.MethodGet, storyID, nil)
	if err != nil {
		return dependachore.Story{}, err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return dependachore.Story{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("string(body) = %+v\n", string(body))
	if err != nil {
		return dependachore.Story{}, err
	}

	story := dependachore.Story{}
	err = json.Unmarshal(body, &story)
	if err != nil {
		return dependachore.Story{}, err
	}

	return story, nil
}

func (c *Client) Chorify(storyID int) error {
	return c.updateStory(storyID, `{"story_type": "chore"}`)
}

func (c *Client) MoveAfter(storyID, afterStoryID int) error {
	return c.updateStory(storyID, fmt.Sprintf(`{"after_id": %d}`, afterStoryID))
}

func (c *Client) createRequest(httpMethod string, storyID int, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s/%d?fields=before_id,after_id,description,story_type,requested_by_id,id", trackerURL, c.storyPath, storyID)
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-TrackerToken", c.apiKey)
	request.Header.Set("Content-Type", "application/json")

	return request, nil
}

func (c *Client) updateStory(storyID int, storyUpdate string) error {
	request, err := c.createRequest(http.MethodPut, storyID, bytes.NewReader([]byte(storyUpdate)))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("story update failed, http status: %d, message: %s", resp.StatusCode, string(respBody))
	}
	return nil
}
