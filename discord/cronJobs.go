package discord

import (
	"fmt"
	"log"

	"github.com/nikojunttila/discord/riot"
	"github.com/nikojunttila/discord/utils"
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
		for i, user := range usersToCheck {
			newMatch, check := riot.CheckLastMatch(user.lastMatch, user.puuID, user.region, apiCfg.apiKey)
			if check {
				usersToCheck[i].lastMatch = newMatch
				result, ranked, err := riot.GetMatch(newMatch, user.puuID, user.region, apiCfg.apiKey)
				if err != nil {
					fmt.Println("Loop check: ", err)
					return
				}
				region2 := "euw1"
				if user.region == "AMERICAS" {
					region2 = "na1"
				}
				if ranked {
					newRankStats, err := riot.RankedStats(user.puuID, apiCfg.apiKey, region2)
					if err != nil {
						fmt.Println(err)
					}
					oldRank := fmt.Sprintf("Old rank %s: %s %dlp\n", user.Tier, user.Rank, user.LeaguePoints)
					newRank := fmt.Sprintf("New rank %s: %s %dlp\n", newRankStats.Tier, newRankStats.Rank, newRankStats.LeaguePoints)
					result += oldRank
					result += newRank
					usersToCheck[i].LeagueEntry = newRankStats
					promote, checker := utils.CheckPromotionDemotion(user.Tier, user.Rank, newRankStats.Tier, newRankStats.Rank)
					if checker {
						result += promote
					}
				}
				sendGameStatus(s, result, channelID)
			}
		}
	})
	c.AddFunc("0 20 * * ?", func() {
		sendCat(s, "660136166515015711")
	})
	/* 	c.AddFunc("@every 5m", func() {
		sendCat(s, "249254722668724225")
	}) */
	c.Start()
}
