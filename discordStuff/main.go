package discord

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/discord/internal/database"

	"github.com/gempir/go-twitch-irc/v4"
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

func Discord() {
	dbURL := utils.GetEnvVariable("DB_URL")
	connection, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cant connect to database", err)
	}
	guildID := utils.GetEnvVariable("GUILD_ID")
	channelID := utils.GetEnvVariable("channel_ID")
	apiKey := utils.GetEnvVariable("riot_API")
	apiCfg = apiConfig{
		DB:     database.New(connection),
		apiKey: apiKey,
	}
	initializeDiscordHandlers()

	//twitch stuff
	twitchCH := utils.GetEnvVariable("tchannel")
	//channel ID that shows messages in discord
	tTodChannel := utils.GetEnvVariable("tTodChannel")
	oauth := utils.GetEnvVariable("oauth")
	client := twitch.NewClient("tpx_bot", oauth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if message.User.DisplayName == "tpx_bot" {
			return
		}
		sendMessageToChannel(s, fmt.Sprintf("%s: %s\n", message.User.DisplayName, message.Message), tTodChannel)
	})

	client.Join(twitchCH)
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		if m.ChannelID == tTodChannel {
			user, _ := s.User(m.Author.ID)
			//m.Author
			response := fmt.Sprintf("%s: %s", user.Username, m.Content)
			client.Say(twitchCH, response)
		}
	})
	go client.Connect()
	defer client.Disconnect()
	//twitch stuff end

	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer s.Close()

	//new guilds need to be registered here for slash commands
	registerDiscordCommands("")
	setupCron(channelID)

	//this creates auto mod rules when uncommented to discord guild
	/* 	ruleID := setupAutoModerationRule(guildID, channelID)
	   	defer s.AutoModerationRuleDelete(guildID, ruleID) */

	utils.WaitForInterruptSignal()
	//true to remove slash commands from bot
	removeCommands := false
	if removeCommands {
		removeRegisteredCommands(guildID)
	}

	log.Println("Gracefully shutting down.")
}

/* func sendTwitchMessage(s *discordgo.Session, m *discordgo.MessageCreate, client *twitch.Client) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID == "400298523263893505" {
		client.Say("randomderppy", m.Content)
	}
} */
