package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type MixpanelEvent struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

func TrackMixpanelEvent(event string, distinctID string, props map[string]interface{}) {
	properties := map[string]interface{}{
		"token":       "6d71d0fef020541a28518a4836bdaa10",
		"distinct_id": distinctID,
		"time":        time.Now().Unix(),
	}

	for k, v := range props {
		properties[k] = v
	}

	payload := MixpanelEvent{
		Event:      event,
		Properties: properties,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Println("Mixpanel JSON marshal error:", err)
		return
	}

	encoded := base64.StdEncoding.EncodeToString(jsonPayload)
	form := url.Values{}
	form.Add("data", encoded)

	log.Printf("[Mixpanel] Sending event: %s for user: %s\n", event, distinctID)
	log.Printf("[Mixpanel] Payload: %s\n", string(jsonPayload))

	resp, err := http.PostForm("https://api.mixpanel.com/track", form)
	if err != nil {
		log.Println("Mixpanel request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Mixpanel tracking failed with status: %s\n", resp.Status)
	}
}
