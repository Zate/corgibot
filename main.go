package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	yaml "gopkg.in/yaml.v2"
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

type apikeys struct {
	BotKey string
}

type corgiAPIresp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (a *apikeys) getAPIKeys(filename string) *apikeys {
	yamlFile, err := ioutil.ReadFile(filename)
	errCheck("", err)
	err = yaml.Unmarshal(yamlFile, a)
	errCheck("", err)
	return a
}

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

	if content == "!corgime" {
		// https://dog.ceo/api/breed/corgi/images/random
		resp, err := http.Get("https://dog.ceo/api/breed/corgi/images/random")
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "Failed to get corgi pic")
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		errCheck("failed to parse response", err)
		corgiPic := corgiAPIresp{}
		err = json.Unmarshal(body, &corgiPic)
		errCheck("unmarshal of json failed", err)

		// resp, err = http.Get(corgiPic.Message)
		// errCheck("failed to download the pic", err)

		// defer resp.Body.Close()

		//pic := io.Reader(resp.Body)

		// pic, err := os.Create("/tmp/corgi.jpg")
		// errCheck("failed to save temp pic", err)

		// _, err = io.Copy(file, resp.Body)
		// errCheck("write failed", err)

		// file.Close

		embed := &discordgo.MessageEmbed{
			Author: &discordgo.MessageEmbedAuthor{},
			Color:  0x00ff00, // Green
			//Description: "This is a discordgo embed",
			// Fields: []*discordgo.MessageEmbedField{
			// 	&discordgo.MessageEmbedField{
			// 		Name:   "I am a field",
			// 		Value:  "I am a value",
			// 		Inline: true,
			// 	},
			// 	&discordgo.MessageEmbedField{
			// 		Name:   "I am a second field",
			// 		Value:  "I am a value",
			// 		Inline: true,
			// 	},
			// },
			Image: &discordgo.MessageEmbedImage{
				URL: corgiPic.Message,
				// },
				// Thumbnail: &discordgo.MessageEmbedThumbnail{
				// 	URL: corgiPic.Message,
			},
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			//Title:     "I am an Embed",
		}

		discord.ChannelMessageSendEmbed(message.ChannelID, embed)

		//discord.ChannelFileSend(message.ChannelID, "Corgi", pic)

	}

	log.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	var k apikeys
	k.getAPIKeys(".secrets.yaml")
	discord, err := discordgo.New("Bot " + k.BotKey)
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
