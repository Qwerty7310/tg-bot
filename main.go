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

	url := os.Getenv("NEIGRY_URL")
	if url == "" {
		log.Fatal("Установите NEIGRY_URL")
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
			resp, err := http.Get(url)
			if err != nil {
				msg := tgbotapi.NewMessage(chatID, "Не удалось сделать запрос!")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Ошибка отправки:", err)
				}
			}

			if resp.StatusCode == http.StatusOK {
				msg := tgbotapi.NewMessage(chatID, "200 OK")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Ошибка отправки:", err)
				}
			} else {
				msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Сайт ответил: %d %s\n", resp.StatusCode, resp.Status))
				_, err := bot.Send(msg)
				if err != nil {
					log.Println("Ошибка отправки:", err)
				}
			}

			resp.Body.Close()
			time.Sleep(15 * time.Second)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
