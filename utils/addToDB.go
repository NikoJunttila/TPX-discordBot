package utils

import (
	"github.com/nikojunttila/discord/riot"
)

type CheckMatch struct {
	name      string
	hashtag   string
	puuID     string
	region    string
	lastMatch string
}

func AddToDB(name string, hashtag string, apiKey string, region string) error {
	puuID, err := riot.GetPuuID(name, hashtag, apiKey)
	if err != nil {
		return err
	}

}
