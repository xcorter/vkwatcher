package telegramclient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/vkclient"
	"log"
	"net/http"
)

type TelegramClient struct {
	tgbotapi *tgbotapi.BotAPI
}

func (t *TelegramClient) SendMessage(
	observable observable.Observable,
	item vkclient.Item,
) {

	text := item.Text
	for _, attach := range item.Attachments {
		if attach.Type == "audio" {
			text += "\n" + attach.Audio.Artist + " - " + attach.Audio.Title
		}
	}

	msg := tgbotapi.NewMessage(observable.ChatId, text)

	_, err := t.tgbotapi.Send(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (t *TelegramClient) GetUpdatesChan() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	return t.tgbotapi.GetUpdatesChan(u)
}

func NewClient(token string, client *http.Client) *TelegramClient {
	bot, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		log.Panic(err)
	}

	//bot.Debug = true
	return &TelegramClient{
		tgbotapi: bot,
	}
}
