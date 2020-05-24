package telegramclient

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/vkclient"
	"log"
	"net/http"
	"strconv"
)

type TelegramClient struct {
	tgbotapi *tgbotapi.BotAPI
}

func (t *TelegramClient) SendMessage(
	observable observable.Observable,
	item vkclient.Item,
) {
	text := ""
	if item.PostType == "post" {
		text = "Новый пост:\n"
	} else if item.PostType == "reply" {
		text = "Трек в комментах:\n"
	}

	text = item.Text + "\n"
	for _, attach := range item.Attachments {
		if attach.Type == "audio" {
			text += attach.Audio.Artist + " - " + attach.Audio.Title + "\n"
		}
	}
	link := t.getLink(item)
	text += link

	fmt.Println(observable.Owner + ":\n" + text)
	msg := tgbotapi.NewMessage(observable.ChatId, text)

	_, err := t.tgbotapi.Send(msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (t *TelegramClient) getLink(item vkclient.Item) string {
	link := "https://vk.com/"
	publicId := item.OwnerId
	if publicId < 0 {
		publicId = publicId * -1
	}
	publicIdStr := strconv.Itoa(publicId)
	postId := strconv.Itoa(item.Id)
	link += "public" + publicIdStr + "?w=wall-" + publicIdStr + "_" + postId
	return link
}

func (t *TelegramClient) SendRawMessage(chatId int64, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)

	_, err := t.tgbotapi.Send(msg)
	if err != nil {
		return err
	}
	return nil
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
