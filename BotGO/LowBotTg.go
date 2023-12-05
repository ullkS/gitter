package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Создание бота,ввод токена
	bot, err := tgbotapi.NewBotAPI("Телеграм токен")
	if err != nil {
		log.Fatal(err)
	}

	// Устанвока режима отладки
	bot.Debug = true

	log.Printf("Авторизация %s", bot.Self.UserName)

	// Обработчик обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Получаем данные от Telegram API
	updates:= bot.GetUpdatesChan(u)

	// Обрабатываем данные
	for update := range updates {
		if update.Message == nil {
			continue
		}
		text := update.Message.Text

		// Отвечаем на сообщение в зависимости от его содержимого
		reply := getReply(text)
		// Отправляем ответное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Функция для получения ответа на сообщение
func getReply(text string) string {
	switch text {
	case "Привет":
		return "Привет!"
	case "Как дела?":
		return "Хорошо, спасибо!"
	case "Что делаешь?":
		return "Отвечаю на сообщения в Telegram :)"
	default:
		return "Извините, я не понимаю вас"
	}
}
