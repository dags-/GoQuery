package goquery

import (
	"net/http"
	"bytes"
	"encoding/json"
)

type Profile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Legacy bool `json:"legacy,omitempty"`
}

const mojangApi = "https://api.mojang.com/profiles/minecraft"

func Profiles(status interface{}) []Profile {
	strings, ok := status.([]string)
	if ok {
		return GetProfiles(strings...)
	}
	return []Profile{}
}

func GetProfiles(names ...string) []Profile {
	payload, payloadErr := json.Marshal(names)
	if payloadErr != nil {
		return []Profile{}
	}

	client := &http.Client{}
	post, postErr := http.NewRequest("POST", mojangApi, bytes.NewBuffer(payload))
	response, postErr := client.Do(post)

	if postErr != nil {
		return []Profile{}
	}

	var profiles []Profile
	decoder := json.NewDecoder(response.Body)
	decodeErr := decoder.Decode(&profiles)

	if decodeErr != nil {
		return []Profile{}
	}
	return profiles
}