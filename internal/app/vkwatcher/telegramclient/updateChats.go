package telegramclient

import (
	"context"
	"fmt"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"strings"
)

type ChatWatcher struct {
	telegramClient *TelegramClient
	provider       *observable.Provider
}

func (c *ChatWatcher) Watch(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop chat watcher")
				return
			default:
				updates, err := c.telegramClient.GetUpdatesChan()

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				for update := range updates {
					chatId := update.Message.Chat.ID
					if update.Message.IsCommand() {
						text := "Привет! Этот бот следит за исполнителем в вк и постит сюда треки, которые были " +
							"запощены. Просто отправьте имя исполнителя и бот начнет слежку.\n Сканирование VK " +
							"проходит с небольшими временными интервалами."
						err := c.telegramClient.SendRawMessage(chatId, text)
						if err != nil {
							fmt.Println(err.Error())
						}
						continue
					}

					amountOfChats := c.provider.GetCountByChatId(chatId)
					if amountOfChats > 1 {
						text := "Достигнут предел испонителей"
						err := c.telegramClient.SendRawMessage(chatId, text)
						if err != nil {
							fmt.Println(err.Error())
						}
						continue
					}

					username := update.Message.Chat.LastName + " " + update.Message.Chat.FirstName + "|" +
						update.Message.Chat.UserName
					artistName := strings.ToLower(update.Message.Text)
					artistName = strings.TrimSpace(artistName)
					ob := observable.NewMusicObservable(username, artistName, chatId)

					c.provider.Save(ob)

					successMessage := "Исполнитель " + artistName + " был успешно добавлен"
					err := c.telegramClient.SendRawMessage(chatId, successMessage)
					if err != nil {
						fmt.Println(err.Error())
					}
					continue
				}
			}
		}
	}()
}

func NewChatWatcher(client *TelegramClient, provider *observable.Provider) ChatWatcher {
	return ChatWatcher{
		telegramClient: client,
		provider:       provider,
	}
}
