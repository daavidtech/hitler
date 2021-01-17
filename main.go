package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

var token string

// var buffer = make([][]byte, 0)

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: hitler -t <bot token>")
		return
	}

	// // Load the sound file.
	// err := loadSound()
	// if err != nil {
	// 	fmt.Println("Error loading sound: ", err)
	// 	fmt.Println("Please copy $GOPATH/src/github.com/bwmarrin/examples/airhorn/airhorn.dca to this directory.")
	// 	return
	// }

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Airhorn is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateStatus(0, "!siegheil")
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "kill yourself") {
		// // Stop speaking
		// vc.Speaking(false)

		// // Sleep for a specificed amount of time before ending.
		// time.Sleep(250 * time.Millisecond)

		// // Disconnect from the provided voice channel.
		// vc.Disconnect()

		panic("Killing my self now")

		return
	}

	if strings.HasPrefix(m.Content, "!siegheil") {
		times := 1

		parts := strings.Split(m.Content, " ")

		var err error

		if len(parts) > 1 {
			times, err = strconv.Atoi(parts[1])

			if err != nil {
				times = 1
			}
		}

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			// Could not find guild.
			return
		}

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				err = playSound(s, g.ID, vs.ChannelID, times)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
}

func loadRandomSound() ([][]byte, error) {
	var buffer = make([][]byte, 0)

	files, _ := ioutil.ReadDir("./quotes")

	numberOfFiles := len(files)

	rnd := rand.Intn(numberOfFiles)

	file := files[rnd]

	fileHandle, _ := os.Open("./quotes/" + file.Name())

	var magic uint32

	binary.Read(fileHandle, binary.LittleEndian, &magic)

	var headerLength int32

	binary.Read(fileHandle, binary.LittleEndian, &headerLength)

	header := make([]byte, headerLength)

	binary.Read(fileHandle, binary.LittleEndian, &header)

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err := binary.Read(fileHandle, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := fileHandle.Close()
			if err != nil {
				return buffer, err
			}
			return buffer, err
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return buffer, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(fileHandle, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}

	fmt.Println("Load file success")

	return buffer, nil
}

// // loadSound attempts to load an encoded sound file from disk.
// func loadSound() error {

// 	file, err := os.Open("airhorn.dca")
// 	if err != nil {
// 		fmt.Println("Error opening dca file :", err)
// 		return err
// 	}

// 	var opuslen int16

// 	for {
// 		// Read opus frame length from dca file.
// 		err = binary.Read(file, binary.LittleEndian, &opuslen)

// 		// If this is the end of the file, just return.
// 		if err == io.EOF || err == io.ErrUnexpectedEOF {
// 			err := file.Close()
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}

// 		if err != nil {
// 			fmt.Println("Error reading from dca file :", err)
// 			return err
// 		}

// 		// Read encoded pcm from dca file.
// 		InBuf := make([]byte, opuslen)
// 		err = binary.Read(file, binary.LittleEndian, &InBuf)

// 		// Should not be any end of file errors
// 		if err != nil {
// 			fmt.Println("Error reading from dca file :", err)
// 			return err
// 		}

// 		// Append encoded pcm data to the buffer.
// 		buffer = append(buffer, InBuf)
// 	}
// }

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, times int) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	for i := 0; i < times; i++ {
		buffer, err := loadRandomSound()

		if err != nil {
			fmt.Println("Error loading file :", err)
		}

		// Send the buffer data.
		for _, buff := range buffer {
			vc.OpusSend <- buff
		}
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
