package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/nikojunttila/discord/riot"
)

var (
	integerOptionMinValue = 1.0

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "hello",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "say hello user",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "user name",
					Required:    true,
				},
			},
		},
		{
			Name:        "highscores",
			Description: "highscores for N-word count",
		},
		{
			Name:        "opgg",
			Description: "Get stats from history max 70",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "acc name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "hashtag",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "europeoramericas",
					Description: "europe or americas",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "number of games",
					MinValue:    &integerOptionMinValue,
					MaxValue:    80,
					Required:    true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
		{
			Name:        "addtofollows",
			Description: "add user to live tracking",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "acc name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "hashtag",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "europeoramericas",
					Description: "europe or americas",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"hello": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			var name string
			if option, ok := optionMap["name"]; ok {
				name = option.StringValue()
			}
			mention := "<@" + name + "> "
			response := fmt.Sprintf("Hello %s!", mention)
			if i.Member.User.ID != "" {
				fmt.Println(i.Member.User.ID)
				if i.Member.User.ID == "249254722668724225" {
					response = "Neekeri"
					sendTag(s, response, i.ChannelID, name)
					return
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
		"highscores": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			response, err := apiCfg.getHighscores()
			if err != nil {
				fmt.Println(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Error getting highscores. Try again later.",
					},
				})
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
				},
			})
		},
		"opgg": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			type infos struct {
				name    string
				hashtag string
				region  string
				amount  int64
			}
			var inf infos

			if option, ok := optionMap["name"]; ok {
				inf.name = option.StringValue()
			}

			if opt, ok := optionMap["hashtag"]; ok {
				inf.hashtag = opt.StringValue()
			}
			if opt, ok := optionMap["europeoramericas"]; ok {
				inf.region = opt.StringValue()
			}
			if opt, ok := optionMap["amount"]; ok {
				inf.amount = opt.IntValue()
			}
			fmt.Println(inf)
			var response string
			puuID, err := riot.GetPuuID(inf.name, inf.hashtag, globalAPi)
			if err != nil {
				fmt.Println("slash puuid get ", err)
				response = fmt.Sprintln(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				})
				return
			}
			history, err := riot.GetMatchHistory(puuID, int(inf.amount), inf.region, globalAPi)
			if err != nil {
				fmt.Println("slash history get ", err)
				response = fmt.Sprintln(err)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				})
				return
			}
			response, err = riot.PrintHistory(history, globalAPi, puuID, inf.region, inf.name)
			if err != nil {
				fmt.Println("slash history print", err)
				response = fmt.Sprintln("slash print history err:", err)
			}
			fmt.Println(response)
			if inf.amount < 15 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: response,
					},
				})
			} else {
				sendGameStatus(s, response, i.ChannelID)
			}
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			//send message hidden
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "Surprise!",
				},
			})
		},
		"addtofollows": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			type infos struct {
				name    string
				hashtag string
				region  string
			}
			var inf infos

			if option, ok := optionMap["name"]; ok {
				inf.name = option.StringValue()
			}

			if opt, ok := optionMap["hashtag"]; ok {
				inf.hashtag = opt.StringValue()
			}
			if opt, ok := optionMap["europeoramericas"]; ok {
				inf.region = opt.StringValue()
			}
			fmt.Println(inf)
			err := AddToDB(inf.name, inf.hashtag, inf.region)
			if err != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Failed to add user to tracking feed. \n",
					},
				})
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Added user %s to tracking feed\n", inf.name),
				},
			})
		},
	}
)
