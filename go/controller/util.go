package controller

import "math/rand"

// character range to choose from for random realm generation
const randomKeyChars = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const randomKeyBaseSize = 19

/**
create two maps<br>
links:	random key -> human readable link<br>
ref:	random key -> target key
*/
func MakeLinks(userName string) (map[string]string, map[string]string) {
	links := make(map[string]string)
	ref := make(map[string]string)
	user := users[userName]
	for _, access := range user.Access {
		if targets[access] != nil {
			key := makeRandomKey()
			for ref[key] != "" { // make sure key is unique
				key = makeRandomKey()
			}
			links[key] = targets[access][0]
			ref[key] = access
		}
	}
	return links, ref
}

func makeRandomKey() string {
	buffer := make([]byte, randomKeyBaseSize)
	for i := range buffer {
		buffer[i] = randomKeyChars[rand.Intn(len(randomKeyChars))]
	}
	return string(buffer)
}

type Refmap struct {
	Refs, Targets map[string]string
}

func NewRefmap(r, t map[string]string) *Refmap {
	return &Refmap{r, t}
}
