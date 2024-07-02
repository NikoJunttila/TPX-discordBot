package discord

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikojunttila/discord/internal/database"
	"github.com/nikojunttila/discord/riot"
)

type checkMatch struct {
	name string
	//hashtag   string
	puuID     string
	region    string
	lastMatch string
	inGame    bool
	riot.LeagueEntry
}

func AddToDB(name string, hashtag string, region string) error {
	puuID, err := riot.GetPuuID(name, hashtag, apiCfg.apiKey)
	if err != nil {
		return err
	}
	dbCtx := context.Background()
	_, err = apiCfg.DB.CreateFollow(dbCtx, database.CreateFollowParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now().UTC(),
		AccountName: name,
		Hashtag:     hashtag,
		Puuid:       puuID,
		Region:      region,
	})
	if err != nil {
		fmt.Println("Error creating follow ", err)
		return err
	}
	return nil
}
func getFollows() ([]checkMatch, error) {
	dbCtx := context.Background()
	var returnUsers []checkMatch
	users, err := apiCfg.DB.GetFollowed(dbCtx)
	if err != nil {
		fmt.Println("Error getting follows ", err)
		return returnUsers, err
	}
	for _, user := range users {
		newUser := checkMatch{
			name:      user.AccountName,
			puuID:     user.Puuid,
			region:    user.Region,
			inGame:    false,
			lastMatch: "", // Initialize lastMatch to avoid nil panic
		}
		lastMatches, err := riot.GetMatchHistory(user.Puuid, 1, user.Region, apiCfg.apiKey)
		if err != nil {
			fmt.Println("Error GetMatchHistory ", err)
			return returnUsers, err
		}
		if len(lastMatches) > 0 {
			newUser.lastMatch = lastMatches[0]
		}
		region2 := "euw1"
		if user.Region == "AMERICAS" {
			region2 = "na1"
		}
		rankedStuff, err := riot.RankedStats(user.Puuid, apiCfg.apiKey, region2)
		if err != nil {
			fmt.Println(err)
			return returnUsers, err
		}
		newUser.LeagueEntry = rankedStuff
		returnUsers = append(returnUsers, newUser)
	}
	return returnUsers, nil
}
