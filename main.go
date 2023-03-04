package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"github.com/YoungPentagonHacker/cryptorate/cryptocurrencyparser"
	"github.com/YoungPentagonHacker/cryptorate/database"
	"github.com/YoungPentagonHacker/cryptorate/rubparser"
	"github.com/YoungPentagonHacker/cryptorate/timemanager"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

var rubCost float64

func main() {

	token, _ := os.LookupEnv("TGBOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			rubCost = rubparser.Parse()
			if rubCost == 0 {
				rubCost = 70
			}
			time.Sleep(5 * 60 * time.Second)
		}

	}()

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}
	var activecommand string
	for update := range updates {
		if update.Message == nil {
			continue
		}
		var msg tgbotapi.MessageConfig

		if !strings.HasPrefix(update.Message.Text, "/") {
			switch activecommand {
			case "timeset":
				err = database.SetTime(update.Message.Chat.ID, update.Message.Text)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong time format!!!")
					bot.Send(msg)
					break
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "The time has been set")
				bot.Send(msg)
			case "addcrypto":
				err = database.AddCrypto(update.Message.Chat.ID, update.Message.Text)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "No cryptocurrency with this name was found, check if it is available on coinmarketcap or send a link")
					bot.Send(msg)
					break
				}
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "The cryptocurrency has been added")
				bot.Send(msg)
			case "sendcourse":
				var resp string
				crypto := update.Message.Text
				if !database.CryptoValidate(crypto) {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Wrong cryptocurrency name")
					bot.Send(msg)
					break
				}
				bucksCost := cryptocurrencyparser.Parse(crypto)
				bucksFloat, _ := strconv.ParseFloat(bucksCost[1:], 64)
				rubFloat := rubCost * bucksFloat
				log.Println(rubCost, bucksCost)

				resp += crypto + " - " + bucksCost + " ₽" + fmt.Sprintf("%.2f", rubFloat) + "\n"
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, resp)
				bot.Send(msg)
			}
		}

		switch update.Message.Command() {
		case "start":
			database.CreateUser(update.Message.Chat.ID)
		case "timeset":
			activecommand = "timeset"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Choose the time at which you would like to receive the price of currencies daily")
			bot.Send(msg)
		case "addcrypto":
			activecommand = "addcrypto"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Send the full name of a crypto or the full link to the coinmarketcap for the cryptocurrency you would like to receive notifications about")
			bot.Send(msg)

		case "launch":
			//if database.GetActive(update.Message.Chat.ID) == true {
			//	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Daily notifications have been already started")
			//	bot.Send(msg)
			//	break
			//}
			msg, err = responseForm(update.Message.Chat.ID)
			if err != nil {
				bot.Send(msg)
				break
			}
			go timemanager.WaitUntil(database.GetTime(update.Message.Chat.ID), func() {
				m, _ := responseForm(update.Message.Chat.ID)
				bot.Send(m)

			}, func() bool { return database.GetOk(update.Message.Chat.ID) })
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Daily notifications started successfully")
			bot.Send(msg)
			//database.SetActive(update.Message.Chat.ID, true)

		case "stop":
			database.SetOk(update.Message.Chat.ID, false)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Daily notifications have been stopped")
			bot.Send(msg)

		case "sendcourse":
			activecommand = "sendcourse"
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Write the name of the crypto you want to receive right now")
			bot.Send(msg)
		}

	}
}

func responseForm(chatid int64) (tgbotapi.MessageConfig, error) {
	crypto := database.GetCrypto(chatid)
	var msg tgbotapi.MessageConfig
	if crypto == nil {
		msg = tgbotapi.NewMessage(chatid, "You must select at least one cryptocurrency!")
		database.SetOk(chatid, false)
		return msg, errors.New("no crypto in db")
	}
	if database.GetTime(chatid) == "" {
		msg = tgbotapi.NewMessage(chatid, "Select the time of the notification!")
		database.SetOk(chatid, false)
		return msg, errors.New("no time in db")

	}
	var resp string
	for i := 0; i < len(crypto); i++ {
		bucksCost := cryptocurrencyparser.Parse(crypto[i])
		bucksFloat, _ := strconv.ParseFloat(bucksCost[1:], 64)
		rubFloat := rubCost * bucksFloat
		log.Println(rubCost, bucksCost)

		resp += crypto[i] + " - " + bucksCost + " ₽" + fmt.Sprintf("%.2f", rubFloat) + "\n"
	}
	msg = tgbotapi.NewMessage(chatid, resp)
	database.SetOk(chatid, true)

	return msg, nil

}
