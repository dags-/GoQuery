package goquery

import (
	"net/http"
	"bytes"
	"encoding/json"
	"sort"
)

type Profile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Legacy bool `json:"legacy,omitempty"`
}

type profileArray []Profile

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

	var profiles profileArray
	decoder := json.NewDecoder(response.Body)
	decodeErr := decoder.Decode(&profiles)

	if decodeErr != nil {
		return []Profile{}
	}
	sort.Sort(profiles)
	return profiles
}

func (profiles profileArray) Len() int {
	return len(profiles)
}

func (profiles profileArray) Less(i, j int) bool {
	return profiles[i].Name < profiles[j].Id
}

func (profiles profileArray) Swap(i, j int)  {
	profiles[i], profiles[j] = profiles[j], profiles[i]
}