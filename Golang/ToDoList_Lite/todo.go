package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tasks := []string{}
	var choice string
	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		fmt.Print("Выберите действие: ")
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			fmt.Print("Введите задачу: ")
			task, _ := reader.ReadString('\n')
			task = task[:len(task)-1] // Удаляем символ новой строки
			tasks = append(tasks, task)
			fmt.Println("Задача добавлена.")
		case "2":
			if len(tasks) == 0 {
				fmt.Println("Список задач пуст.")
			} else {
				fmt.Println("Список задач:")
				for i, task := range tasks {
					fmt.Printf("%d. %s\n", i+1, task)
				}
			}
		case "3":
			fmt.Println("Программа завершена.")
			return
		default:
			fmt.Println("Неверный выбор. Попробуйте еще раз.")
		}
	}
}

func printMenu() {
	fmt.Println("\nМеню:")
	fmt.Println("1. Добавить задачу")
	fmt.Println("2. Показать список задач")
	fmt.Println("3. Выйти")
}
