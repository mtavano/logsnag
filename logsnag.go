package logsnag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type LogSnag struct {
	Token   string
	Project string
}

func (logsnag *LogSnag) GetProject() string {
	return logsnag.Project
}

type PublishRequest struct {
	Project     string            `json:"project,omitempty"`
	Channel     string            `json:"channel,omitempty"`
	Event       string            `json:"event,omitempty"`
	Description string            `json:"description,omitempty"`
	Icon        string            `json:"icon,omitempty"`
	Notify      bool              `json:"notify,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
	Parser      string            `json:"parser,omitempty"`
	UserID      string            `json:"user_id,omitempty"`
	Timestamp   int64             `json:"timestamp,omitempty"`
}

func (logsnag *LogSnag) Publish(input *PublishRequest) error {
	if input.Channel == "" || input.Event == "" {
		return errors.New("logsnag: LogSnag.Publish missing one of required fields project, channel, or event")
	}
	input.Project = logsnag.GetProject()

	baseURL := "https://api.logsnag.com/v1/log"

	body, err := json.Marshal(input)
	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish json.Marshal error")
	}

	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewReader(body))

	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish http.NewRequest error")
	}

	req.Header.Add("Authorization", "Bearer "+logsnag.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "logsnag: LogSnag.Publish http.DefaultClient.Do error")
	}

	if res.StatusCode != http.StatusOK || res.StatusCode != http.StatusCreated {
		return fmt.Errorf("logsnag: LogSnag.Publish unexpected http response status <%d> error", res.StatusCode)
	}

	return nil
}

// TODO: mofify me when need insight
func (logsnag *LogSnag) Insight(title string, value string, icon string) bool {
	url := "https://api.logsnag.com/v1/insight"
	method := "POST"

	payload := strings.NewReader(`{
		"project": "` + logsnag.GetProject() + `",
		"title": "` + title + `",
		"value": "` + value + `",
		"icon": "` + icon + `"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Add("Authorization", "Bearer "+logsnag.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	_, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func NewLogSnag(token string, project string) LogSnag {
	return LogSnag{
		Token:   token,
		Project: project,
	}
}

func main() {
	logSnag := NewLogSnag(
		"d67d3443e793dad29d9c94df76838367",
		"ferry-times",
	)

	logSnag.Insight(
		"waitlist",    // Channel
		"User Joined", // Event
		"🛥️",          // Icon
	)
}
