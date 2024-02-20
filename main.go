package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/discord/riot"
	"github.com/nikojunttila/discord/utils"
	"github.com/robfig/cron/v3"
)

var s *discordgo.Session
var globalAPi string

func init() {
	godotenv.Load()
	BotToken := os.Getenv("DISCORD")
	if BotToken == "" {
		fmt.Println("port not found in env")
	}
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
	godotenv.Load()
	BotToken := os.Getenv("DISCORD")
	if BotToken == "" {
		fmt.Println("port not found in env")
	}
	//guildID := os.Getenv("GUILD_ID")
	guildID := os.Getenv("TPX_ID")
	if guildID == "" {
		fmt.Println("port not found in env")
	}
	//channelID := os.Getenv("channel_ID")
	channelID := os.Getenv("general2")
	if channelID == "" {
		fmt.Println("port not found in env")
	}
	apiKey := os.Getenv("riot_API")
	if apiKey == "" {
		fmt.Println("port not found in env")
	}
	globalAPi = apiKey
	s.Identify.Intents |= discordgo.IntentAutoModerationExecution
	s.Identify.Intents |= discordgo.IntentMessageContent
	s.Identify.Intents |= discordgo.IntentsGuildMessages
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		fmt.Println("registered command")
		registeredCommands[i] = cmd
	}

	defer s.Close()

	c := cron.New()
	alphaChecker, puuID := riot.InitStats("Alphass", "EUW", "EUROPE", apiKey)
	bziChecker, puuID2 := riot.InitStats("Best Voli Iraq", "EUW", "EUROPE", apiKey)
	lenkkisChecker, puuID3 := riot.InitStats("lenkkis", "SNEED", "EUROPE", apiKey)
	kadeemChecker, puuID4 := riot.InitStats("kadeem", "718", "AMERICAS", apiKey)
	var check bool
	c.AddFunc("@every 1m", func() {
		alphaChecker, check = riot.CheckLastMatch(alphaChecker, puuID, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(alphaChecker, puuID, "EUROPE", apiKey)
			sendGameStatus(s, result, channelID)
		}
		bziChecker, check = riot.CheckLastMatch(bziChecker, puuID2, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(bziChecker, puuID2, "EUROPE", apiKey)
			sendGameStatus(s, result, channelID)
		}
		lenkkisChecker, check = riot.CheckLastMatch(lenkkisChecker, puuID3, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(lenkkisChecker, puuID3, "EUROPE", apiKey)
			sendGameStatus(s, result, channelID)
		}
		kadeemChecker, check = riot.CheckLastMatch(kadeemChecker, puuID4, "AMERICAS", apiKey)
		if check {
			result, _ := riot.GetMatch(kadeemChecker, puuID4, "AMERICAS", apiKey)
			sendGameStatus(s, result, channelID)
		}
	})
	c.AddFunc("@every 50m", func() {
		randomNumber := rand.Intn(6)
		if randomNumber == 3 {
			responses := []string{"top gap", "neekeri", "java enjoyer", "If you were any more inbred, you'd be a sandwich", "Your map awareness is so bad, even Twisted Fate wouldn't ult to save you.", "Not even Olaf ult could prevent you from being disabled",
				"I'd call you cancer but at least cancer gets kills", "If i wanted to kill myself i'd jump up to your ego and jump down to your IQ.", "Even the mars curiosity rover has faster reaction time than you", "Even Christopher Columbus had better map awareness than you"}
			rand2 := rand.Intn(len(responses) + 1)
			sendTag(s, responses[rand2], channelID, "685511498641965089")
		}

	})
	c.AddFunc("@every 1d", func() {
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
	})
	c.Start()
	enabled := true
	rule, err := s.AutoModerationRuleCreate(guildID, &discordgo.AutoModerationRule{
		Name:        "NNZ",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*nigger*", "neekeri", "ngr", "nigga"},
			/* 	RegexPatterns: []string{}, */
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
	defer s.AutoModerationRuleDelete(guildID, rule.ID)
	s.AddHandler(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		nWordCalc, err := utils.IncrementAndWriteToFile("nWordCount.txt")
		if err != nil {
			fmt.Println("Error nwordCounter:", err)
		}
		s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Hei! Ei N-pommia tässä discordissa! N-sana sanottu %d kertaa\n", nWordCalc))
	})
	s.AddHandler(messageCreate)
	// Wait here until CTRL-C or other term signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	removeCommands := true
	if removeCommands {
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
	log.Println("Gracefully shutting down.")
}
