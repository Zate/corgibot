package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func errCheck(msg string, err error) {
	if err != nil {
		log.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

var (
	commandPrefix string
	botID         string
)

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	content := message.Content

	if content == "!test" {
		discord.ChannelMessageSend(message.ChannelID, "Testing..")
	}

	log.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	discord, err := discordgo.New("Bot NDYzOTEzMjkzNzk5NjIwNjA4.Dh53rQ.9JHgyd4W80wCiVjxyFu7akph3WI")
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Admiring Corgi Butt")
		if err != nil {
			log.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		log.Printf("CorgiBot has started on %d servers", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}
