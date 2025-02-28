package logsnag

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

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

func (logsnag *LogSnag) PublishLegacy(
	channel string,
	event string,
	icon string,
	tags map[string]any,
	notify bool,
) bool {
	url := "https://api.logsnag.com/v1/log"
	method := "POST"

	// Create the description string from map
	var pairs []string
	for key, value := range tags {
		pairs = append(pairs, fmt.Sprintf(`%s: %v`, key, value))
	}

	description := strings.Join(pairs, ", ")

	rawPayload := `{
		"project": "` + logsnag.GetProject() + `",
		"channel": "` + channel + `",
		"event": "` + event + `",
		"description": "` + description + `",
		"icon": "` + icon + `",
		"notify": "` + strconv.FormatBool(notify) + `"
	}`
	fmt.Println("--------")
	fmt.Println(rawPayload)
	fmt.Println("--------")

	payload := strings.NewReader(rawPayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}

	bearerToken := "Bearer " + logsnag.Token
	fmt.Printf("TOKEN: <%s>\n", bearerToken)
	req.Header.Add("Authorization", bearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
