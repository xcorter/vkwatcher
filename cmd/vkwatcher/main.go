package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/observable"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/telegramclient"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/vkclient"
	"github.com/xcorter/vkwatcher/internal/app/vkwatcher/watcher"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("run")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := &http.Client{}

	db := getDB()
	provider := observable.NewProvider(db)
	manager := observable.NewManager(provider)

	vkclientapi := vkclient.NewClient(client, os.Getenv("VK_API_KEY"))
	telegramClient := getTelegramClient()
	w := watcher.NewWatcher(
		manager,
		provider,
		vkclientapi,
		telegramClient,
	)
	ctx, cancel := context.WithCancel(context.Background())

	chatWatcher := telegramclient.NewChatWatcher(telegramClient, provider)
	chatWatcher.Watch(ctx)

	//Регистрируем сигналы SIGINT и SIGTERM, чтобы завершить работу красиво
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		fmt.Println("Нажали crtl+c. Полылаем сингал завершения")
		cancel()
		fmt.Println("Ждем завершения вокреков")
	}()

	w.Start(ctx)
}

func getTelegramClient() *telegramclient.TelegramClient {
	client := getTelegramHTTPClient()

	return telegramclient.NewClient(
		os.Getenv("TELEGRAM_API_KEY"),
		client,
	)
}

func getTelegramHTTPClient() *http.Client {
	//creating the proxyURL
	proxyStr := os.Getenv("USE_PROXY")
	if proxyStr == "" {
		return &http.Client{
			//Timeout: 10 * time.Second,
		}
	}
	proxyURL, err := url2.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}
	//adding the proxy settings to the Transport object
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	//adding the Transport object to the http Client
	return &http.Client{
		Transport: transport,
		//Timeout: 5 * time.Second,
	}
}

func getDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./observable.db")
	checkErr(err)
	tableSql := "CREATE TABLE IF NOT EXISTS observable (owner VARCHAR(255), value VARCHAR(255), type INT, last_scan INT, chat_id VARCHAR(255))"
	_, err = db.Exec(tableSql)
	if err != nil {
		panic(err)
	}
	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
