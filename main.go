package main

import (
	"context"
	trnsl "discord_bot/src/translator"
	w "discord_bot/src/weather"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

var commandList = map[string]string{
	"!help":           "Get descriptions of available commands.",
	"!weather":        "Get city weather information by city name.\nSample: !weather Almaty",
	"!translate":      "Translate text to target language.\nSample: !translate ru Hello world!\n!translate russian Hello world!",
	"!translate_lang": "Show available languages for translate.",
}

// loadEnv load env for project
func loadEnv() {
	// loads DB settings from .env into the system
	if err := godotenv.Load("./bot.env"); err != nil {
		log.Print("ERROR. No .env file found.", err)
	}
}

// getCommandDescr return info about bot commands
func getCommandDescr() *discordgo.MessageSend {
	var fields []*discordgo.MessageEmbedField

	for key, val := range commandList {
		fields = append(fields, &discordgo.MessageEmbedField{Name: key, Value: val, Inline: false})
	}
	data := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:   discordgo.EmbedTypeRich,
				Title:  "Commands description:",
				Fields: fields,
			},
		},
	}
	return data
}

// startBot start discord bot
func startBot() {
	botToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Print("ERROR. Create bot err.", err)
	}
	discord.AddHandler(mainHandler)
	discord.Open()
	defer discord.Close()

	log.Print("Bot running....")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// mainHandler catch user command for exec
func mainHandler(discord *discordgo.Session, message *discordgo.MessageCreate) {
	ctx := context.Background()

	//ignore bot messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	switch {
	case strings.Contains(message.Content, "!help"):
		helpInfo := getCommandDescr()
		discord.ChannelMessageSendComplex(message.ChannelID, helpInfo)

	case strings.Contains(message.Content, "!translate_lang"):
		res := trnsl.GetSupportedLang(ctx)
		discord.ChannelMessageSend(message.ChannelID, res)

	case strings.Contains(message.Content, "!translate"):
		re := regexp.MustCompile("!translate (?P<lang>[A-Za-zА-яа-я]+) (?P<text>.+)")
		res := re.FindStringSubmatch(message.Content)
		if len(res) > 2 {
			langTo := res[re.SubexpIndex("lang")]
			trnsl.IsSupportedLang(ctx, langTo)
			text := res[re.SubexpIndex("text")]
			trnsl.IsSupportedLang(ctx, langTo)
			translation, err := trnsl.Translate(ctx, text, langTo)
			if err != nil {
				discord.ChannelMessageSend(message.ChannelID, err.Error())
			} else {
				discord.ChannelMessageSend(message.ChannelID, translation)
			}
		} else {
			discord.ChannelMessageSend(message.ChannelID, "Сommand used incorrectly!")
		}

	case strings.Contains(message.Content, "!weather"):
		city, ok := strings.CutPrefix(message.Content, "!weather ")
		if !ok {
			discord.ChannelMessageSend(message.ChannelID, "Сommand used incorrectly!")
		} else {
			if currentWeather, err := w.GetByLocName(city); err != nil {
				discord.ChannelMessageSend(message.ChannelID, err.Error())
			} else {
				discord.ChannelMessageSendComplex(message.ChannelID, currentWeather)
			}
		}
	}
}

func main() {
	loadEnv()
	startBot()
}
