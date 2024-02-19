package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nikojunttila/discord/riot"
	"github.com/nikojunttila/discord/utils"
	"github.com/robfig/cron/v3"
)

func main() {
	godotenv.Load()
	Token := os.Getenv("DISCORD")
	if Token == "" {
		fmt.Println("port not found in env")
	}
	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		fmt.Println("port not found in env")
	}

	channelID := os.Getenv("channel_ID")
	if channelID == "" {
		fmt.Println("port not found in env")
	}
	apiKey := os.Getenv("riot_API")
	if apiKey == "" {
		fmt.Println("port not found in env")
	}
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents |= discordgo.IntentAutoModerationExecution
	dg.Identify.Intents |= discordgo.IntentMessageContent
	dg.Identify.Intents |= discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	c := cron.New()
	alphaChecker, puuID := riot.InitStats("Alphass", "EUW", "EUROPE", apiKey)
	bziChecker, puuID2 := riot.InitStats("Best Voli Iraq", "EUW", "EUROPE", apiKey)
	lenkkisChecker, puuID3 := riot.InitStats("lenkkis", "SNEED", "EUROPE", apiKey)
	kadeemChecker, puuID4 := riot.InitStats("kadeem", "718", "AMERICAS", apiKey)
	var check bool
	c.AddFunc("@every 2m", func() {
		alphaChecker, check = riot.CheckLastMatch(alphaChecker, puuID, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(alphaChecker, puuID, "EUROPE", apiKey)
			sendGameStatus(dg, result, channelID)
		}
		bziChecker, check = riot.CheckLastMatch(bziChecker, puuID2, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(bziChecker, puuID2, "EUROPE", apiKey)
			sendGameStatus(dg, result, channelID)
		}
		lenkkisChecker, check = riot.CheckLastMatch(lenkkisChecker, puuID3, "EUROPE", apiKey)
		if check {
			result, _ := riot.GetMatch(lenkkisChecker, puuID3, "EUROPE", apiKey)
			sendGameStatus(dg, result, channelID)
		}
		kadeemChecker, check = riot.CheckLastMatch(kadeemChecker, puuID4, "AMERICAS", apiKey)
		if check {
			result, _ := riot.GetMatch(kadeemChecker, puuID4, "AMERICAS", apiKey)
			sendGameStatus(dg, result, channelID)
		}
	})
	c.AddFunc("@every 1d", func() {
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
	})
	c.Start()
	enabled := true
	rule, err := dg.AutoModerationRuleCreate(guildID, &discordgo.AutoModerationRule{
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
	defer dg.AutoModerationRuleDelete(guildID, rule.ID)
	dg.AddHandler(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		/* 	nWordCalc, err := utils.IncrementAndWriteToFile("nWordCount.txt") */
		nWordCalc, err := utils.IncrementAndWriteToFile()
		if err != nil {
			fmt.Println("Error:", err)
		}
		s.ChannelMessageSend(e.ChannelID, fmt.Sprintf("Hei! Ei N-pommia tässä discordissa! N-sana sanottu %d kertaa\n", nWordCalc))
	})
	dg.AddHandler(messageCreate)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
