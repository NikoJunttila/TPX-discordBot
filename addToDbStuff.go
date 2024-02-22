package main

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
}

func AddToDB(name string, hashtag string, region string) error {
	puuID, err := riot.GetPuuID(name, hashtag, globalAPi)
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

/*
	 func getFollows() ([]checkMatch, error) {
		dbCtx := context.Background()
		var returnUsers []checkMatch
		users, err := apiCfg.DB.GetFollowed(dbCtx)
		if err != nil {
			fmt.Println("Error getting follows ", err)
			return returnUsers, err
		}
		for i, user := range users {
			returnUsers[i].name = user.AccountName
			returnUsers[i].puuID = user.Puuid
			returnUsers[i].region = user.Region
			lastMatches, err := riot.GetMatchHistory(user.Puuid, 1, user.Region, globalAPi)
			if err != nil {
				fmt.Println("Error GetMatchHistory ", err)
				return returnUsers, err
			}
			returnUsers[i].lastMatch = lastMatches[0]
		}
		return returnUsers, nil
	}
*/
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
			lastMatch: "", // Initialize lastMatch to avoid nil panic
		}
		lastMatches, err := riot.GetMatchHistory(user.Puuid, 1, user.Region, globalAPi)
		if err != nil {
			fmt.Println("Error GetMatchHistory ", err)
			return returnUsers, err
		}
		if len(lastMatches) > 0 {
			newUser.lastMatch = lastMatches[0]
		}
		returnUsers = append(returnUsers, newUser)
	}
	return returnUsers, nil
}
