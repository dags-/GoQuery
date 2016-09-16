package profiles

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const mojangApi = "https://api.mojang.com/profiles/minecraft"

type Profile struct {
	Id   string
	Name string
}

type Profiles []Profile

func LookupProfiles(players []string) Profiles {
	return LookupProfile(players...)
}

func LookupProfile(players ...string) Profiles {
	payload, payloadErr := json.Marshal(players)
	if payloadErr != nil {
		return Profiles{}
	}

	client := &http.Client{}
	post, postErr := http.NewRequest("POST", mojangApi, bytes.NewBuffer(payload))
	response, postErr := client.Do(post)

	if postErr != nil {
		return Profiles{}
	}

	var profiles Profiles
	decoder := json.NewDecoder(response.Body)
	decodeErr := decoder.Decode(&profiles)

	if decodeErr != nil {
		return Profiles{}
	}

	return profiles
}

func (profiles *Profiles) ToJson() string {
	return toJson(profiles, false)
}

func (profiles *Profiles) ToPrettyJson() string {
	return toJson(profiles, true)
}

func toJson(input interface{}, pretty bool) string {
	var prefix, indent = "", ""
	if pretty {
		indent = "    "
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(input)
	if err != nil {
		return ""
	}
	return string(buffer.Bytes())
}
