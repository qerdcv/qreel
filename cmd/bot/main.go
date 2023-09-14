package main

import (
	"bytes"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/qerdcv/qreel/pkg/reelser"
)

// Bot written for personal purpose
func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalln(err.Error())
	}

	r := reelser.New()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "download":
			reelURL := update.Message.CommandArguments()
			if reelURL == "" {
				if _, rErr := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You forgot to provide url!\nPlease use next syntax /download <url>")); rErr != nil {
					log.Println("ERROR: send message", rErr.Error())
					continue
				}
				continue
			}

			log.Println("INFO: Downloading video", reelURL)

			url, rErr := r.GetVideoURL(reelURL)
			if rErr != nil {
				log.Println("ERROR: get video url", rErr.Error())
				continue
			}

			if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Downloading!")); err != nil {
				log.Println("ERROR: send message", err.Error())
				continue
			}

			buf := new(bytes.Buffer)
			video := tgbotapi.NewVideo(update.Message.Chat.ID, tgbotapi.FileReader{
				Name:   "video.mp4",
				Reader: buf,
			})

			if rErr = r.DownloadReel(url, buf); rErr != nil {
				log.Println("ERROR: download reel", rErr.Error())
				continue
			}

			username := update.Message.From.UserName
			if username == "" {
				username = update.Message.From.FirstName + update.Message.From.LastName
			}

			video.Caption = "Reels from " + username
			if _, err = bot.Send(video); err != nil {
				log.Println("ERROR: send", err.Error())
				continue
			}

			log.Println("INFO: done with", reelURL)

			continue
		}

		if _, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command :(")); err != nil {
			log.Println("ERROR: send", err.Error())
		}
	}
}
