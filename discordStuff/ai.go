package discord

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/nikojunttila/discord/ai"
	"github.com/nikojunttila/discord/utils"
)

const (
	// aiMaxHistory caps how many past messages are kept per user so the context
	// sent to the model (and the token cost) stays bounded.
	aiMaxHistory   = 5
	aiSystemPrompt = "You are a helpful assistant in a Discord server. Answer concisely and helpfully."
	// discordMsgLimit is Discord's maximum message content length.
	discordMsgLimit = 2000
)

// aiHistory holds a rolling, per-user conversation history. Handlers run in
// their own goroutine (discordgo dispatches with SyncEvents=false), so all
// access is guarded by aiHistoryMu.
var (
	aiHistory   = make(map[string][]ai.Message)
	aiHistoryMu sync.Mutex
)

func askHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	var question string
	if option, ok := optionMap["question"]; ok {
		question = option.StringValue()
	}

	// Acknowledge immediately: the AI call takes several seconds and Discord
	// requires a response within 3 seconds. We edit this deferred message once
	// the answer is ready.
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		fmt.Println("ask defer:", err)
		return
	}

	userID := interactionUserID(i)

	// Copy the user's history so the network call happens without holding the lock.
	aiHistoryMu.Lock()
	history := append([]ai.Message(nil), aiHistory[userID]...)
	aiHistoryMu.Unlock()

	messages := make([]ai.Message, 0, len(history)+2)
	messages = append(messages, ai.Message{Role: "system", Content: aiSystemPrompt})
	messages = append(messages, history...)
	messages = append(messages, ai.Message{Role: "user", Content: question})

	answer, err := ai.Ask(apiCfg.kimiKey, messages)
	if err != nil {
		fmt.Println("ask kimi:", err)
		sendInteractionText(s, i, "Sorry, I couldn't get an answer right now. Try again later.")
		return
	}

	// Record the exchange and trim to the most recent aiMaxHistory messages.
	aiHistoryMu.Lock()
	updated := append(
		aiHistory[userID],
		ai.Message{Role: "user", Content: question},
		ai.Message{Role: "assistant", Content: answer},
	)
	if len(updated) > aiMaxHistory {
		updated = updated[len(updated)-aiMaxHistory:]
	}
	aiHistory[userID] = updated
	aiHistoryMu.Unlock()

	sendInteractionText(s, i, answer)
}

func aiClearHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	userID := interactionUserID(i)
	aiHistoryMu.Lock()
	delete(aiHistory, userID)
	aiHistoryMu.Unlock()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Cleared your AI conversation history.",
		},
	})
}

// interactionUserID returns the invoking user's ID, whether the command was run
// in a guild (i.Member) or a DM (i.User).
func interactionUserID(i *discordgo.InteractionCreate) string {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID
	}
	if i.User != nil {
		return i.User.ID
	}
	return ""
}

// sendInteractionText delivers text as the reply to a deferred interaction,
// splitting it across multiple messages to respect Discord's 2000-char limit.
func sendInteractionText(s *discordgo.Session, i *discordgo.InteractionCreate, text string) {
	chunks := utils.SplitMessage(text, discordMsgLimit)
	if len(chunks) == 0 {
		chunks = []string{"(no response)"}
	}

	// The first chunk edits the original (deferred) response.
	content := chunks[0]
	if _, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &content}); err != nil {
		fmt.Println("ask edit response:", err)
		return
	}

	// Any remaining chunks are sent as public followups.
	for _, chunk := range chunks[1:] {
		if _, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{Content: chunk}); err != nil {
			fmt.Println("ask followup:", err)
		}
	}
}
