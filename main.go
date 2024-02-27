package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/discord/internal/database"
	"github.com/nikojunttila/discord/riot"

	_ "github.com/lib/pq"
	"github.com/nikojunttila/discord/utils"
	"github.com/robfig/cron/v3"
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

	registerDiscordCommands(guildID)
	c := cron.New()
	usersToCheck, err := getFollows()
	if err != nil {
		log.Panic(err)
	}

	c.AddFunc("@every 1m", func() {
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
	c.AddFunc("@every 1d", func() {
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
	})
	c.Start()

	ruleID := setupAutoModerationRule(guildID, channelID)
	defer s.AutoModerationRuleDelete(guildID, ruleID)

	// Wait here until CTRL-C or other term signal is received.
	utils.WaitForInterruptSignal()

	removeCommands := true
	if removeCommands {
		removeRegisteredCommands(guildID)
	}

	log.Println("Gracefully shutting down.")
}
func registerDiscordCommands(guildID string) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func initializeDiscordHandlers() {
	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates | discordgo.IntentMessageContent | discordgo.IntentAutoModerationExecution
	s.AddHandler(ready)
	s.AddHandler(apiCfg.badWordCounter)
	s.AddHandler(messageCreate)

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
}

func setupAutoModerationRule(guildID, channelID string) string {
	enabled := true
	rule, err := s.AutoModerationRuleCreate(guildID, &discordgo.AutoModerationRule{
		Name:        "NNZ",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: utils.Automod,
		},
		Enabled: &enabled,
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{
				ChannelID: channelID,
			}},
		},
	})
	if err != nil {
		panic(err)
	}
	return rule.ID
}

func removeRegisteredCommands(guildID string) {
	log.Println("Removing commands...")
	registeredCommands2, err := s.ApplicationCommands(s.State.User.ID, guildID)
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands2 {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
