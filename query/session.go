package goquery

import (
	"net/http"
	"encoding/json"
	"encoding/base64"
)

const mojangSessionApi = "https://sessionserver.mojang.com/session/minecraft/profile/"

type Session struct {
	Profile
	Skin string `json:"skin,omitempty"`
}

type mojangSession struct {
	Profile

	Properties []struct {
		Name  string
		Value string
	}
}

type mojangSkin struct {
	Timestamp   int64
	ProfileId   string
	ProfileName string
	Textures    struct {
			    SKIN struct {
					 Url string
				 }
		    }
}

func Sessions(profiles ...Profile) []Session {
	sessions := make([]Session, len(profiles))
	for i := range profiles {
		sessions[i] = GetSession(profiles[i].Id)
	}
	return sessions
}

func GetSession(id string) Session {
	session := getSession(id)
	skin := decodeSkinUrl(session)
	return Session{Profile{session.Id, session.Name, session.Legacy}, skin}
}

func getSession(id string) mojangSession {
	url := mojangSessionApi + id
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		return mojangSession{}
	}

	var sess = mojangSession{}

	sessionDecoder := json.NewDecoder(response.Body)
	sessionErr := sessionDecoder.Decode(&sess)

	if sessionErr != nil {
		return mojangSession{}
	}
	return sess
}

func decodeSkinUrl(sess mojangSession) string {
	if len(sess.Properties) == 0 {
		return ""
	}

	value := sess.Properties[0].Value
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}

	var skn mojangSkin
	jsErr := json.Unmarshal(decoded, &skn)
	if jsErr != nil {
		return ""
	}

	return skn.Textures.SKIN.Url
}