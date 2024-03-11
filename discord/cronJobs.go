package discord

import (
	"fmt"
	"log"

	"github.com/nikojunttila/discord/riot"
	"github.com/robfig/cron/v3"
)

func setupCron(channelID string) {
	usersToCheck, err := getFollows()
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println(usersToCheck)
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		//fmt.Println("checking for games...")
		for i, user := range usersToCheck {
			newMatch, check := riot.CheckLastMatch(user.lastMatch, user.puuID, user.region, apiCfg.apiKey)
			if check {
				usersToCheck[i].lastMatch = newMatch
				result, err := riot.GetMatch(newMatch, user.puuID, user.region, apiCfg.apiKey)
				if err != nil {
					fmt.Println("Loop check: ", err)
					return
				}
				sendGameStatus(s, result, channelID)
			}
		}
	})
	/* 	c.AddFunc("@every 6h", func() {
		sendCat(s, "660136166515015711")
	}) */
	c.AddFunc("0 20 * * ?", func() {
		sendCat(s, "660136166515015711")
	})
	/* 	c.AddFunc("@every 5m", func() {
		sendCat(s, "249254722668724225")
	}) */
	c.Start()
}
