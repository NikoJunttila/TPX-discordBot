package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/discord/internal/database"

	_ "github.com/lib/pq"
	"github.com/nikojunttila/discord/utils"
)

type apiConfig struct {
	DB     *database.Queries
	apiKey string
}

var (
	s      *discordgo.Session
	apiCfg apiConfig
)

func init() {
	godotenv.Load()
	BotToken := utils.GetEnvVariable("DISCORD")
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	dbURL := utils.GetEnvVariable("DB_URL")
	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cant connect to database", err)
	}
	//guildID := utils.GetEnvVariable("GUILD_ID")
	guildID := utils.GetEnvVariable("TPX_ID")
	//channelID := utils.GetEnvVariable("channel_ID")
	channelID := utils.GetEnvVariable("general2")
	apiKey := utils.GetEnvVariable("riot_API")
	apiCfg = apiConfig{
		DB:     database.New(connection),
		apiKey: apiKey,
	}

	initializeDiscordHandlers()
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer s.Close()

	//new guilds need to be registered here for slash commands
	registerDiscordCommands("")
	setupCron(channelID)

	//auto mod rules
	/* 	ruleID := setupAutoModerationRule(guildID, channelID)
	   	defer s.AutoModerationRuleDelete(guildID, ruleID) */

	// Wait here until CTRL-C or other term signal is received.
	utils.WaitForInterruptSignal()
	//true to remove slash commands from bot
	removeCommands := false
	if removeCommands {
		removeRegisteredCommands(guildID)
	}

	log.Println("Gracefully shutting down.")
}
