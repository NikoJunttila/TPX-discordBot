package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/nikojunttila/discord/riot"
	"github.com/nikojunttila/discord/utils"
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
			Name: "live",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "check players live game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "name",
					Description: "account name",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "hashtag",
					Description: "String option",
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
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "amount",
					Description: "number of games",
					MinValue:    &integerOptionMinValue,
					MaxValue:    80,
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "region",
					Description: "leave empty for eu",
					Required:    false,
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
					Name:        "region",
					Description: "leave empty for eu",
					Required:    false,
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
				if i.Member.User.ID == "249254722668724225" {
					response = utils.InsultRes()
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
		"live": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}
			var name string
			var hashtag string
			if option, ok := optionMap["name"]; ok {
				name = option.StringValue()
			}
			if option, ok := optionMap["hashtag"]; ok {
				hashtag = option.StringValue()
			}
			response, err := riot.LiveGamePlayersStatsFormattedToString(apiCfg.apiKey, name, hashtag)
			if err != nil {
				fmt.Println(err)
				return
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
			if opt, ok := optionMap["amount"]; ok {
				inf.amount = opt.IntValue()
			}
			if opt, ok := optionMap["region"]; ok {
				value := strings.ToLower(opt.StringValue())
				switch value {
				case "eune", "euw", "ru", "tr":
					inf.region = "EUROPE"
				case "na", "br":
					inf.region = "AMERICAS"
				case "kr", "jp":
					inf.region = "ASIA"
				default:
					inf.region = "EUROPE"
				}
			} else {
				inf.region = "EUROPE"
			}
			fmt.Println(inf)
			var response string
			puuID, err := riot.GetPuuID(inf.name, inf.hashtag, apiCfg.apiKey)
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
			history, err := riot.GetMatchHistory(puuID, int(inf.amount), inf.region, apiCfg.apiKey)
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
			response, err = riot.PrintHistory(history, apiCfg.apiKey, puuID, inf.region, inf.name)
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
				sendMessageToChannel(s, response, i.ChannelID)
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

			if opt, ok := optionMap["region"]; ok {
				value := strings.ToLower(opt.StringValue())
				switch value {
				case "eune", "euw", "ru", "tr":
					inf.region = "EUROPE"
				case "na", "br":
					inf.region = "AMERICAS"
				case "kr", "jp":
					inf.region = "ASIA"
				default:
					inf.region = "EUROPE"
				}
			} else {
				inf.region = "EUROPE"
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
