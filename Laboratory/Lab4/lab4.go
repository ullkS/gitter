package main

import (
	"fmt"
	"sync"
	"time"
)

const bufferSize = 20

type Buffer struct {  //Кольцевой буфер
	data       [bufferSize]int
	readIndex  int
	writeIndex int
	count      int
	mutex      sync.Mutex
	notEmpty   *sync.Cond
	notFull    *sync.Cond
}

func NewBuffer() *Buffer {
	buffer := &Buffer{}
	buffer.notEmpty = sync.NewCond(&buffer.mutex)
	buffer.notFull = sync.NewCond(&buffer.mutex)
	return buffer
}

func (b *Buffer) Write(value int, priority int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for b.count == bufferSize {
		b.notFull.Wait()
	}

	b.data[b.writeIndex] = value
	b.writeIndex = (b.writeIndex + 1) % bufferSize
	b.count++

	fmt.Printf("Писатель записал сообщение %d. Имеет приоритет %d\n", value, priority)

	b.notEmpty.Broadcast()
}

func (b *Buffer) Read(priority int) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for b.count == 0 {
		b.notEmpty.Wait()
	}

	value := b.data[b.readIndex]
	b.readIndex = (b.readIndex + 1) % bufferSize
	b.count--

	fmt.Printf("Читатель прочитал сообщение %d. Приоритет %d\n", value, priority)

	b.notFull.Broadcast()

	return value
}

func main() {
	buffer := NewBuffer()

	var wg sync.WaitGroup

	// Создаем писателей. Возможно изменение количество писателей, сопоставив читателей и буфер.
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		start:=time.Now()
		go func(writerID int) {
			CauntMessage:=0
			for j := 1; j <= 10; j++ {
				buffer.Write(j, writerID)
				CauntMessage++
			}
			fmt.Println("Писатель:",writerID,"Количество сообщений:",CauntMessage,"Время:",time.Now().Sub(start))
			wg.Done()
		}(i)
	
	}

	// Создаем читателей. Возможно изменение количество читателей, сопоставив писателей и буфер.
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		start:=time.Now()
		go func(readerID int) {
			CauntMessage:= 0 
			for j := 1; j <= 7; j++ {
				buffer.Read(readerID)
				CauntMessage++
			}
			fmt.Println("Читатель:",readerID,"Количество сообщений:",CauntMessage,"Время:",time.Now().Sub(start))
			wg.Done()
		}(i)
	}

	wg.Wait()
}