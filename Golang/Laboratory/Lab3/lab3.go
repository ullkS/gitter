package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

const bufferSize = 1 
const exitKeyword = "exit"

var buffer = make(chan string, bufferSize)
var wg sync.WaitGroup

func main() {
	writerCount := 5
	readerCount := 3

	wg.Add(writerCount + readerCount + 1)

	defer fmt.Println("Все сообщения выведены")
	fmt.Printf("Введите сообщение (для выходда введите %s): ", exitKeyword)

	for i := 0; i < writerCount; i++ {
		go writer(i)
	}

	for i := 0; i < readerCount; i++ {
		go reader(i)
	}

	wg.Wait()
}

func writer(id int) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		message := scanner.Text()
		if message == exitKeyword {
			close(buffer)
		}
		buffer <- message
		fmt.Printf("Писатель %d записал сообщение: %s\n", id, message)
	}
}

func reader(id int) {
	for { 
		message := <-buffer
		start := time.Now()
		if message == exitKeyword {
			break
		}
		
		fmt.Printf("Читатель %d прочитал сообщение: %s\n", id, message)
		end := time.Now()
		fmt.Println("Время чтения/записи:", start.Sub(end))
	}
	wg.Done()
}
