package main

import (
	"io/ioutil"
	"log"

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
	botID string
)

const (
	svrName         = "Ark Servers Bot"
	svrStatus       = "Monitoring Ark Servers"
	adminBit  int64 = 0x8
)

type apikeys struct {
	BotKey string
}

func getAPIKeys(filename string) string {
	yamlFile, err := ioutil.ReadFile(filename)
	log.Printf("%v", string(yamlFile))
	errCheck("", err)
	b := &apikeys{}
	err = yaml.Unmarshal(yamlFile, b)
	errCheck("", err)
	log.Printf("%+v", b)
	return b.BotKey
}

func commandHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == botID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	channel, err := discord.State.Channel(message.ChannelID)
	errCheck("Error getting info on "+message.ChannelID, err)
	gid := channel.GuildID

	// errCheck("Error getting roles for "+gid, err)
	// adminRID := ""
	// for _, v := range roles {
	// 	// need to check for 0x00000008 against user flags.

	// 	if v.Permissions&adminBit == 1 {
	// 		adminRID = v.ID
	// 	}
	// }
	member, err := discord.GuildMember(gid, message.Author.ID)
	errCheck("Error getting roles for "+message.Author.Username, err)
	// admin := false
	roles, err := discord.GuildRoles(gid)
	errCheck("Error getting roles for "+gid, err)
	for _, v := range roles {
		for _, r := range member.Roles { 
			if v.Name == r {
				if v.Permissions&adminBit == 1 {
					// admin = true
				}

			}
		}
	}
	content := message.Content

	if content == "!test" {
		discord.ChannelMessageSend(message.ChannelID, "Testing..")
		// errCheck(""+msg.ChannelID, err)
		log.Printf("Command: %+v Message: %+v || From: %s\n", content, message.Message, message.Author)
	}

	// if content == "!corgime" {
	// 	// https://dog.ceo/api/breed/corgi/images/random
	// 	// resp, err := http.Get("https://dog.ceo/api/breed/corgi/images/random")
	// 	// if err != nil {
	// 	// 	discord.ChannelMessageSend(message.ChannelID, "Failed to get corgi pic")
	// 	// 	return
	// 	// }
	// 	// body, err := ioutil.ReadAll(resp.Body)
	// 	// defer resp.Body.Close()
	// 	// errCheck("failed to parse response", err)
	// 	// corgiPic := corgiAPIresp{}
	// 	// err = json.Unmarshal(body, &corgiPic)
	// 	// errCheck("unmarshal of json failed", err)

	// 	corgi := randomCorgi()

	// 	embed := &discordgo.MessageEmbed{
	// 		Author: &discordgo.MessageEmbedAuthor{},
	// 		Color:  0x9542f4, // Green
	// 		Image: &discordgo.MessageEmbedImage{
	// 			//URL: corgiPic.Message,
	// 			URL: corgi,
	// 		},
	// 		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
	// 	}
	// 	discord.ChannelMessageSendEmbed(message.ChannelID, embed)
	// 	log.Printf("Command: %+v Message: %+v || From: %s\n", content, message.Message, message.Author)
	// }

}

// func randomCorgi() string {
// 	url := "https://www.google.com/search?q=corgi&tbm=isch"

// 	client := &http.Client{}

// 	req, err := http.NewRequest("GET", url, nil)
// 	errCheck("New Request", err)

// 	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; CrOS x86_64 12607.34.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.42 Safari/537.36")

// 	resp, err := client.Do(req)
// 	errCheck("client do", err)

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)

// 	errCheck("io read body", err)
// 	// ou":"http://spectrum-sitecore-spectrumbrands.netdna-ssl.com/~/media/Pet/Furminator/Images/Solution%20Center%20Images/Feature%20Images/corgi.jpg"
// 	re := regexp.MustCompile("ou\":\"(http[^\"]+)\"")
// 	//re := regexp.MustCompile("src=\"(http[^\"]+)\"")
// 	matches := re.FindAllStringSubmatch(string(body), -1)
// 	log.Println(string(body))
// 	corgis := make([]string, len(matches))
// 	log.Println(len(matches))
// 	for index, match := range matches {
// 		log.Println(match[0])
// 		corgis[index] = match[0]
// 	}

// 	//seed with nanoseconds to get make sure unique random number returned
// 	rand.Seed(time.Now().UnixNano())

// 	corgi := corgis[rand.Intn(len(corgis))]
// 	//get random image url and print to stdout
// 	log.Println(corgi)
// 	return corgi

// }

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	k := new(apikeys)
	k.BotKey = getAPIKeys(".secrets.yaml")
	// botkey, err := ioutil.ReadFile("/run/secrets/botkey")
	// errCheck("Not able to read botkey secret", err)
	// k.BotKey = string(botkey)
	discord, err := discordgo.New("Bot " + k.BotKey)
	log.Println(k.BotKey)
	log.Println(discord.User("@me"))
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)
	botID = user.ID
	discord.AddHandler(commandHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		// err = discord.
		// UpdateStatus(0, svrStatus)
		// if err != nil {
		// 	log.Println("Error attempting to set my status")
		// }
		servers := discord.State.Guilds
		log.Printf(svrName+" has started on %d servers", len(servers))
	})

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)

	//chans, err := discord.GuildChannels(discord.)
	defer discord.Close()

	//commandPrefix := "!"

	<-make(chan struct{})

}

// embed := &discordgo.MessageEmbed{
// 	Author: &discordgo.MessageEmbedAuthor{},
// 	Color:  0x9542f4, // Green
// 	Description: "This is a discordgo embed",
// 	Fields: []*discordgo.MessageEmbedField{
// 		&discordgo.MessageEmbedField{
// 			Name:   "I am a field",
// 			Value:  "I am a value",
// 			Inline: true,
// 		},
// 		&discordgo.MessageEmbedField{
// 			Name:   "I am a second field",
// 			Value:  "I am a value",
// 			Inline: true,
// 		},
// 	},
// 	Image: &discordgo.MessageEmbedImage{
// 		URL: corgiPic.Message,
// 		},
// 		Thumbnail: &discordgo.MessageEmbedThumbnail{
// 			URL: corgiPic.Message,
// 	},
// 	Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
// 	Title:     "I am an Embed",
// }
