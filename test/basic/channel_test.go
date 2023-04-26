package basic

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestChan(t *testing.T) {
	stringsChan := make(chan chan string, 100)
	defer close(stringsChan)
	var wg sync.WaitGroup
	chanGroup := make(chan struct{}, 20)
	defer close(chanGroup)
	go func() {
		for chanItem := range stringsChan {
			go func(chanItem chan string) {
				chanGroup <- struct{}{}
				i, _ := strconv.Atoi(<-chanItem)
				time.Sleep(time.Second)
				chanItem <- strconv.Itoa(i) + "\t" + time.Now().String()
				<-chanGroup
			}(chanItem)
		}
	}()
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(idx int) {
			itemChan := make(chan string, 0)
			stringsChan <- itemChan
			itemChan <- strconv.Itoa(idx)
			fmt.Println(<-itemChan)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
