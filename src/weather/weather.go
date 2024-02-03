package weather

import (
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strconv"
	"sync"
)

var (
	owmKey string //API key for OpenWeatherMap
	once   sync.Once
)

// GetByLocName get weather info from OpenWeatherMap API by location name
func GetByLocName(name string) (*discordgo.MessageSend, error) {
	once.Do(func() { owmKey = os.Getenv("OWM_KEY") }) // once read key from env
	weather, err := owm.NewCurrent("C", "RU", owmKey) // celsius with Russian output
	if err != nil {
		log.Print("ERROR. OpenWeatherMap init error.", err)
		return nil, fmt.Errorf("OpenWeatherMap init error")
	}
	if err := weather.CurrentByName(name); err != nil {
		log.Print("ERROR. OpenWeatherMap get weather error.", err)
		return &discordgo.MessageSend{}, fmt.Errorf("OpenWeatherMap get weather error")
	}

	// read weather info
	city := weather.Name
	if len(weather.Weather) == 0 {
		log.Print("ERROR. OpenWeatherMap choose incorrect location.", err)
		return &discordgo.MessageSend{}, fmt.Errorf("choose incorrect location")
	}
	conditions := weather.Weather[0].Description
	temperature := strconv.FormatFloat(weather.Main.Temp, 'f', 2, 64)
	humidity := strconv.Itoa(weather.Main.Humidity)
	wind := strconv.FormatFloat(weather.Wind.Speed, 'f', 2, 64)

	//result message
	data := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:  discordgo.EmbedTypeRich,
			Title: "Погода в " + city,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Условия",
					Value:  conditions,
					Inline: false,
				},
				{
					Name:   "Температура",
					Value:  temperature + "°C",
					Inline: false,
				},
				{
					Name:   "Влажность",
					Value:  humidity + "%",
					Inline: false,
				},
				{
					Name:   "Ветер",
					Value:  wind + " м/c",
					Inline: false,
				},
			},
		},
		},
	}
	return data, nil

}
