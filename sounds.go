package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var buffer = make([][]byte, 0)

func loadSound(sound string) error {
	fileToPlay := fmt.Sprint("sounds/", sound, ".dca")
	fmt.Println(fileToPlay)
	file, err := os.Open(fileToPlay)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16
	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func playSound(s *discordgo.Session, guildID, channelID, sound string) (err error) {
	log.Println("Playing sound: ", sound)
	err = loadSound(sound)
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		return
	}
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
	vc.Speaking(false)
	time.Sleep(250 * time.Millisecond)
	vc.Disconnect()
	buffer = make([][]byte, 0)
	return nil
}

var (
	previousVoiceStates = make(map[string]*discordgo.VoiceState)
	mutex               = &sync.Mutex{}
)

type usersSounds struct {
	id    string
	sound string
}

func voiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	var users = []usersSounds{
		{id: "685511498641965089", sound: "bzi"},
		{id: "249254722668724225", sound: "uwu"},
		{id: "223070624438943745", sound: "fart"},
	}
	mutex.Lock()
	defer mutex.Unlock()
	//hard coded guild id swap later
	guild, _ := s.State.Guild("615649589621686272")
	for _, vs := range guild.VoiceStates {
		user, _ := s.User(vs.UserID)
		fmt.Print(user.Username, ", ")
		if previousVoiceState, ok := previousVoiceStates[vs.UserID]; ok {
			if previousVoiceState.ChannelID != vs.ChannelID {
				fmt.Printf("User %s (%s#%s) has moved from channel %s to %s\n", vs.UserID, user.Username, user.Discriminator, previousVoiceState.ChannelID, vs.ChannelID)
			}
		} else {
			log.Printf("User %s (%s#%s) has joined channel %s\n", vs.UserID, user.Username, user.Discriminator, vs.ChannelID)
			for _, u := range users {
				if vs.UserID == u.id {
					//hard coded guild id swap later
					err := playSound(s, "615649589621686272", vs.ChannelID, u.sound)
					if err != nil {
						fmt.Println("Error playing sound:", err)
					}
				}
			}
		}
		previousVoiceStates[vs.UserID] = vs
	}
	fmt.Println("")
	time.Sleep(3 * time.Second)
}
