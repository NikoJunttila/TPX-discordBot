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
			if !user.inGame {
				_, err := riot.LiveGamePlayersPuuIDSv5(apiCfg.apiKey, user.puuID)
				if err == nil {
					liveGame, err := riot.LiveGamePlayersStatsPuuIDSkipToString(apiCfg.apiKey, user.puuID, user.name)
					if err == nil {
						usersToCheck[i].inGame = true
						sendMessageToChannel(s, liveGame, channelID)
					}
				}
			}
			//match checking here before live game check
			newMatch, check := riot.CheckLastMatch(user.lastMatch, user.puuID, user.region, apiCfg.apiKey)
			if check {
				usersToCheck[i].lastMatch = newMatch
				usersToCheck[i].inGame = false
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
					oldRank := fmt.Sprintf("%s: %s %dlp --> ", user.Tier, user.Rank, user.LeaguePoints)
					newRank := fmt.Sprintf("%s: %s %dlp\n", newRankStats.Tier, newRankStats.Rank, newRankStats.LeaguePoints)
					result += oldRank
					result += newRank
					usersToCheck[i].LeagueEntry = newRankStats
					promote, checker := utils.CheckPromotionDemotion(user.Tier, user.Rank, newRankStats.Tier, newRankStats.Rank)
					if checker {
						result += promote
					}
				}
				sendMessageToChannel(s, result, channelID)
			}
		}
	})
	c.AddFunc("0 16 * * SUN", func() {
		sendCat(s, "660136166515015711")
	})
	c.AddFunc("0 8-20 * * *", func() {
		sendCatLottery(s, "660136166515015711")
	})
	c.Start()
}
