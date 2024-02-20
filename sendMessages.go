package main

import (
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func sendGameStatus(s *discordgo.Session, m string, ch string) {
	if len(m) < 2 {
		return
	}
	s.ChannelMessageSend(ch, m)
}
func sendTag(s *discordgo.Session, m string, ch string, userID string) {
	if len(m) < 2 {
		return
	}

	mention := "<@" + userID + "> " + m

	_, err := s.ChannelMessageSend(ch, mention)
	if err != nil {
		fmt.Println("err mentioning bzi", err)
	}
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.ID == "685511498641965089" {
		randomNumber := rand.Intn(11)
		fmt.Println(randomNumber)
		if randomNumber == 5 {
			responses := []string{"top gap", "neekeri", "java enjoyer", "If you were any more inbred, you'd be a sandwich", "Your map awareness is so bad, even Twisted Fate wouldn't ult to save you.", "Not even Olaf ult could prevent you from being disabled",
				"I'd call you cancer but at least cancer gets kills", "If i wanted to kill myself i'd jump up to your ego and jump down to your IQ.", "Even the mars curiosity rover has faster reaction time than you", "Even Christopher Columbus had better map awareness than you"}
			rand2 := rand.Intn(len(responses) + 1)
			//s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@bzi %s\n", responses[rand2]))
			s.ChannelMessageSendReply(m.ChannelID, responses[rand2], m.Reference())
		}
	}
	if m.Content == "!hello" {
		// Reply to the user
		_, err := s.ChannelMessageSendReply(m.ChannelID, "Hello, I'm your friendly bot!", m.Reference())
		if err != nil {
			fmt.Println("Error sending reply:", err)
		}
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
	// In this example, we only care about messages that are "ping".
	if m.Content != "!Ping" {
		return
	}

	// We create the private channel with the user who sent the message.
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}
	// Then we send the message through the channel we created.
	_, err = s.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}
