package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	width  = 20
	height = 10
)

type point struct {
	x, y int
}

type snake struct {
	head    point
	body    []point
	growth  int
	dir     string
	running bool
}

var (
	food  point
	score int
)

func main() {
	snake := snake{
		head:    point{width / 2, height / 2},
		body:    []point{},
		growth:  0,
		dir:     "right",
		running: true,
	}
	
	food = generateFood()

	// Вывод инструкции
	fmt.Println("Добро пожаловать в игру Змейка!")
	fmt.Println("Используйте WASD для управления. Нажмите 'q' для выхода.")

	// Запускаем игровой цикл
	for snake.running {
		draw(snake)
		processInput(&snake)
		update(&snake)
		time.Sleep(100 * time.Millisecond)
	}
}

// Очистка экрана
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// Генерация случайной пищи
func generateFood() point {
	rand.Seed(time.Now().UnixNano())
	return point{rand.Intn(width), rand.Intn(height)}
}

// Отрисовка игрового поля
func draw(s snake) {
	clearScreen()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == s.head.x && y == s.head.y {
				fmt.Print("O")
			} else if contains(s.body, point{x, y}) {
				fmt.Print("o")
			} else if x == food.x && y == food.y {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}

	fmt.Println("Счет:", score)
}

// Проверка, содержится ли точка в теле змеи
func contains(body []point, p point) bool {
	for _, b := range body {
		if b == p {
			return true
		}
	}
	return false
}

// Обработка пользовательского ввода
func processInput(s *snake) {
	var input string
	fmt.Scanln(&input)

	switch input {
	case "w":
		if s.dir != "down" {
			s.dir = "up"
		}
	case "s":
		if s.dir != "up" {
			s.dir = "down"
		}
	case "a":
		if s.dir != "right" {
			s.dir = "left"
		}
	case "d":
		if s.dir != "left" {
			s.dir = "right"
		}
	case "q":
		s.running = false
	}
}

// Обновление состояния игры
func update(s *snake) {
	// Движение змеи
	switch s.dir {
	case "up":
		s.head.y--
	case "down":
		s.head.y++
	case "left":
		s.head.x--
	case "right":
		s.head.x++
	}

	// Проверка столкновения с границей
	if s.head.x < 0 || s.head.x >= width || s.head.y < 0 || s.head.y >= height {
		s.running = false
		return
	}

	// Проверка столкновения с телом
	if contains(s.body, s.head) {
		s.running = false
		return
	}

	// Проверка столкновения с пищей
	if s.head == food {
		score++
		s.growth++
		food = generateFood()
	}

	// Увеличение длины змеи
	if s.growth > 0 {
		s.body = append([]point{s.head}, s.body...)
		s.growth--
	} else {
		// Смещение тела змеи
		if len(s.body) > 0 {
			copy(s.body[1:], s.body)
			s.body[0] = s.head
		}
	}
}
