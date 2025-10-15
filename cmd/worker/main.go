package main

import (
	"fmt"
	"sync"
	"time"
)

// job — работа одной итерации
func job(i int) {
	// тут ваша логика

	fmt.Println("do job", i)
	time.Sleep(time.Duration(1) * time.Second)

}

func main() {
	const parallel = 2
	const N = 10 // сколько итераций всего

	sem := make(chan struct{}, parallel) // семафор на 2 слота
	var wg sync.WaitGroup
	start := time.Now()
	for i := 1; i < N; i++ {
		sem <- struct{}{} // займём слот
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			defer func() {
				println("goroutine exit", i)
				<-sem
			}() // освободим слот
			job(i)
		}(i)

	}
	wg.Wait() // все задания выполнены
	println(start.Format("2006-01-02 15:04:05"))
	println(time.Now().Format("2006-01-02 15:04:05"))

}
