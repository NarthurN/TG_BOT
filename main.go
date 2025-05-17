package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 👉 Задание 1: Создай структуру TeamScore!
// Подсказка: нам нужна мапа для хранения очков каждого участника
type TeamScore struct {
	scores map[string]int
}

// 👉 Задание 2: Напиши функцию создания нового счетчика
func NewTeamScore() *TeamScore {
	return &TeamScore{scores: make(map[string]int)}
}

// 👉 Задание 3: Напиши функцию добавления очков
func (ts *TeamScore) AddScore(name string) {
	ts.scores[name]++
}

// 👉 Задание 4: Напиши функцию получения всех очков
func (ts *TeamScore) GetScores() string {
	if len(ts.scores) == 0 {
		return "Пока нет результатов. Добавьте участников командой /add [имя]"
	}

	// // Создаем слайс для сортировки
	// type player struct {
	// 	name  string
	// 	score int
	// }
	// var players []player

	// // Заполняем слайс
	// for name, score := range ts.scores {
	// 	players = append(players, player{name, score})
	// }

	var players []struct {
		name  string
		score int
	}

	for name, score := range ts.scores {
		players = append(players, struct {
			name  string
			score int
		}{name, score})
	}

	// Сортируем по убыванию очков
	sort.Slice(players, func(i, j int) bool {
		if players[i].score == players[j].score {
			return players[i].name < players[j].name
		}
		return players[i].score > players[j].score
	})

	var result strings.Builder
	result.WriteString("🏆 Топ игроков:\n\n")
	// Добавляем эмодзи для первых трех мест
	medals := []string{"🥇", "🥈", "🥉"}

	for i, p := range players {
		if i < len(medals) {
			result.WriteString(fmt.Sprintf("%s %-15s: %d очков\n", medals[i], p.name, p.score))
		} else {
			result.WriteString(fmt.Sprintf("🔹 %-15s: %d очков\n", p.name, p.score))
		}
	}
	return result.String()
}

func main() {
	// 👉 Задание 5: Инициализируй бота с твоим токеном
	token := os.Getenv("TGApi")
	if token == "" {
		log.Panic("TGApi не установлен")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic("Ошибка при создании бота:", err)
	}

	// Создаем счетчик
	teamScore := NewTeamScore()

	// Настраиваем получение обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	fmt.Println("🤖 Бот запущен и ждет команд!")

	// 👉 Задание 6: Напиши обработку команд
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")

		switch update.Message.Command() {
		case "start":
			msg.Text = "Привет! Я бот для учета очков.\n\n" +
				"Доступные команды:\n" +
				"/add [имя] - добавить очко участнику\n" +
				"/score - показать текущие результаты\n" +
				"/help - показать это сообщение\n" +
				"/newGame - новая игра"
		case "help":
			msg.Text = "Команды\n" +
				"Доступные команды:\n" +
				"/add [имя] - добавить очко участнику\n" +
				"/score - показать текущие результаты\n" +
				"/help - показать это сообщение\n" +
				"/newGame - новая игра"
		case "score":
			msg.Text = teamScore.GetScores()
		case "add":
			name := strings.TrimSpace(update.Message.CommandArguments())
			if name == "" {
				msg.Text = "Используйте: /add [имя]\nПример: /add Алексей"
			} else {
				teamScore.AddScore(name)
				msg.Text = fmt.Sprintf("✅ %s получил 1 очко!", name)
			}
		case "newGame":
			clear(teamScore.scores)
			msg.Text = "Новая игра началась!"
		default:
			if update.Message.IsCommand() {
				msg.Text = "Неизвестная команда. Введите /help для списка команд"
			} else {
				continue
			}
		}

		// Отправляем сообщение
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
