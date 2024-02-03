package main

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"strings"
)

var (
	botToken string
	owmKey   string
)

func loadEnv() {
	// loads DB settings from .env into the system
	if err := godotenv.Load("./bot.env"); err != nil {
		log.Print("No .env file found.", err)
	}
	botToken = os.Getenv("TOKEN")
	owmKey = os.Getenv("OWM_KEY")
}

func Run() {
	// create a session
	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Print("Create bot err.", err)
	}

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	discord.Open()
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}
	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "weather"):
		discord.ChannelMessageSend(message.ChannelID, "I can help with that! Use '!zip <zip code>'")
	case strings.Contains(message.Content, "bot"):
		discord.ChannelMessageSend(message.ChannelID, "Hi there!")

	}

}

// https://api.openweathermap.org/data/3.0/onecall?lat=33.44&lon=-94.04&exclude=hourly,daily&appid=72289e1ca528512ed0f2ac0400c0b9d3
func main() {
	loadEnv()

	w, err := owm.NewCurrent("F", "ru", owmKey) // fahrenheit (imperial) with Russian output
	if err != nil {
		log.Fatalln(err)
	}
	w.CurrentByName("Almaty")
	fmt.Println(w)
	Run()
	fmt.Println(botToken)

}
