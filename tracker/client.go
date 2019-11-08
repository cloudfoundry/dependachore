package tracker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	apiKey    string
	projectID int
}

type Story struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	RequestedByID int    `json:"requested_by_id"`
	StoryType     string `json:"story_type"`
	BeforeID      int    `json:"before_id"`
	AfterID       int    `json:"after_id"`
	ID            int    `json:"id"`
}

func NewClient(apiKey string, projectID int) *Client {
	return &Client{
		apiKey:    apiKey,
		projectID: projectID,
	}
}

func (c *Client) Get(storyID int) (Story, error) {
	resp, err := c.makeRequest(http.MethodGet, c.storyUrl(storyID,
		"id",
		"name",
		"description",
		"story_type",
		"before_id",
		"after_id",
		"requested_by_id",
	), nil)
	if err != nil {
		return Story{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Story{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Story{}, fmt.Errorf("getting story failed, http status: %d, message: %s", resp.StatusCode, string(body))
	}

	story := Story{}
	err = json.Unmarshal(body, &story)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}

func (c *Client) CreateFeature(name, description string) (Story, error) {
	return c.createStory(name, description, "feature")
}

func (c *Client) CreateRelease(name, description string) (Story, error) {
	return c.createStory(name, description, "release")
}

func (c *Client) createStory(name, description, storyType string) (Story, error) {
	reqBody := fmt.Sprintf(`{"name": "%s", "description": "%s", "story_type": "%s"}`, name, description, storyType)
	resp, err := c.makeRequest(http.MethodPost, c.storiesUrl(), strings.NewReader(reqBody))
	if err != nil {
		return Story{}, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Story{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Story{}, fmt.Errorf("creating story failed, http status: %d, message: %s", resp.StatusCode, string(respBody))
	}

	story := Story{}
	err = json.Unmarshal(respBody, &story)
	if err != nil {
		return Story{}, err
	}

	return story, nil
}

func (c *Client) Delete(storyID int) error {
	resp, err := c.makeRequest(http.MethodDelete, c.storyUrl(storyID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("creating story failed, http status: %d, message: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func (c *Client) MoveAndChorify(storyID, afterStoryID int) error {
	storyUpdate := fmt.Sprintf(`{"story_type": "chore", "after_id": %d}`, afterStoryID)
	resp, err := c.makeRequest(http.MethodPut, c.storyUrl(storyID), bytes.NewReader([]byte(storyUpdate)))
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

func (c *Client) makeRequest(httpMethod, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-TrackerToken", c.apiKey)
	request.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(request)
}

func (c *Client) storiesUrl() string {
	return fmt.Sprintf("https://www.pivotaltracker.com/services/v5/projects/%d/stories", c.projectID)
}

func (c *Client) storyUrl(storyID int, fields ...string) string {
	url := fmt.Sprintf("%s/%d", c.storiesUrl(), storyID)
	if len(fields) > 0 {
		url = fmt.Sprintf("%s?fields=%s", url, strings.Join(fields, ","))
	}
	return url
}
