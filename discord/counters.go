package discord

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nikojunttila/discord/internal/database"
)

// Define your custom handler function
func (apiCfg *apiConfig) badWordCounter(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
	wordCount := 1
	guildCount := 1
	dbCtx := context.Background()
	user, err := apiCfg.DB.GetUser(dbCtx, e.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		apiCfg.DB.CreateUser(dbCtx, database.CreateUserParams{
			ID:        e.UserID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Count:     1,
		})
	} else if err != nil {
		fmt.Println(err)
		return
	} else {
		if user.Count > 0 {
			apiCfg.DB.UpdateUser(dbCtx, database.UpdateUserParams{
				ID:    e.UserID,
				Count: 1,
			})
			wordCount = int(user.Count)
			wordCount++
		}
	}
	guild, err := apiCfg.DB.GetGuild(dbCtx, e.GuildID)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("New guild")
		apiCfg.DB.CreateGuild(dbCtx, database.CreateGuildParams{
			ID:        e.GuildID,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Count:     1,
		})
	} else if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("update guild:", guild.Count)
		if guild.Count > 0 {
			apiCfg.DB.UpdateGuild(dbCtx, database.UpdateGuildParams{
				ID:    e.GuildID,
				Count: 1,
			})
			guildCount = int(guild.Count)
			guildCount++
		}
	}
	mention := "<@" + e.UserID + "> "
	message := fmt.Sprintf("Hei! Ei N-pommia tässä discordissa! %s on sanonut sen sanan %d kertaa.\n Yhteensä sanottu tässä killassa %d kertaa!", mention, wordCount, guildCount)
	//fmt.Println(count, err)
	s.ChannelMessageSend(e.ChannelID, message)
}
func (api *apiConfig) getHighscores() (string, error) {
	dbCtx := context.Background()
	highscores, err := api.DB.HighscoreUsers(dbCtx)
	if err != nil {
		fmt.Println("Getting highscores:", err)
		return "", err
	}
	var highscoresMessage string
	highscoresMessage = "Highscores for N-WORD count: \n"
	for i, user := range highscores {
		userName, err := s.User(user.ID)
		if err != nil {
			fmt.Println("Error fetching user information:", err)
			return "", err
		}
		if len(userName.Username) > 0 {
			highscoresMessage += fmt.Sprintf("%d. %s - %d\n", i+1, userName.Username, user.Count)
		} else {
			highscoresMessage += fmt.Sprintf("%d. Unknown - %d\n", i+1, user.Count)
		}
	}
	return highscoresMessage, nil
}
