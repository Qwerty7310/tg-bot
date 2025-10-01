package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("Установите TELEGRAM_TOKEN")
	}

	url := []string{
		"https://comedyconcert.ru/event/neigry-december25",
		"https://comedyconcert.ru/event/neigry-decembe25",
		"https://comedyconcert.ru/event/neigry-decemb25",
		"https://comedyconcert.ru/event/neigry-decem25",
		"https://comedyconcert.ru/event/neigry-dece25",
		"https://comedyconcert.ru/event/neigry-dec25",
		"https://comedyconcert.ru/event/neigry-de25",
		"https://comedyconcert.ru/event/neigry-d25",
		"https://comedyconcert.ru/event/neigry-november25",
		"https://comedyconcert.ru/event/neigry-novembe25",
		"https://comedyconcert.ru/event/neigry-novemb25",
		"https://comedyconcert.ru/event/neigry-novem25",
		"https://comedyconcert.ru/event/neigry-nove25",
		"https://comedyconcert.ru/event/neigry-nov25",
		"https://comedyconcert.ru/event/neigry-no25",
		"https://comedyconcert.ru/event/neigry-n25",
		"https://comedyconcert.ru/event/neigry-sept25",
	}

	chatID := int64(1622492999)

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Настраиваем получение апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обрабатываем входящие сообщения
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// простой эхо-бот
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Ошибка отправки:", err)
			}
		}
	}()

	go func() {
		for {
			for _, u := range url {
				resp, err := http.Get(u)
				if err != nil {
					log.Println("Не удалось сделать запрос:", err)
					continue
				}

				if resp.StatusCode == http.StatusOK {
					msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("200 OK\n%s", u))
					_, err := bot.Send(msg)
					if err != nil {
						log.Println("Ошибка отправки:", err)
					}
				}

				resp.Body.Close()
			}
			time.Sleep(30 * time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
