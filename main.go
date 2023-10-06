package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	employee_count = 5
	start_time     = "2022-01-01 15:00:00"
)

func main() {
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(employee_count)

	var meat []string
	//add beef, pork, chicken to slice
	AddMeat(&meat)

	//shuffle slice to make it random
	ShuffleMeat(&meat)

	var meat_channel chan string = make(chan string, 30)
	defer close(meat_channel)
	//add all meat to channel
	for _, m := range meat {
		meat_channel <- m
	}

	//process start
	var names []string = []string{"A", "B", "C", "D", "E"}
	for i := 0; i < employee_count; i++ {
		go Process(names[i], meat_channel, wg)
	}

	wg.Wait()
}

func AddMeat(meat *[]string) {
	for i := 0; i < 10; i++ {
		*meat = append(*meat, "beef")
	}
	for i := 0; i < 7; i++ {
		*meat = append(*meat, "pork")
	}
	for i := 0; i < 5; i++ {
		*meat = append(*meat, "chicken")
	}
}

func ShuffleMeat(meat *[]string) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(*meat), func(i int, j int) {
		(*meat)[i], (*meat)[j] = (*meat)[j], (*meat)[i]
	})
}

func Process(name string, meat_channel <-chan string, wg *sync.WaitGroup) {
	var process_time int
	var curr_time time.Time
	curr_time, _ = time.Parse(time.DateTime, start_time)
	defer wg.Done()

FOR:
	for {
		select {

		//continously get meat from channel
		case m := <-meat_channel:
			fmt.Printf("%s在 %s取得%s\n", name, curr_time.Format(time.DateTime), m)

			//get sleep time depends on type of meat
			switch {
			case m == "beef":
				process_time = 1
			case m == "pork":
				process_time = 2
			case m == "chicken":
				process_time = 3
			}

			//sleep and add time
			time.Sleep(time.Duration(process_time) * time.Second)
			curr_time = curr_time.Add(time.Duration(process_time) * time.Second)

			fmt.Printf("%s在 %s處理完%s\n", name, curr_time.Format(time.DateTime), m)

		//if channel is empty, end goroutine
		default:
			break FOR
		}
	}
}
