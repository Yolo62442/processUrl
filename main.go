package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var urlsChan chan string

func main() {
	workProcNum := runtime.NumCPU()
	runtime.GOMAXPROCS(workProcNum)

	urlsChan = make(chan string, 10)

	scanner := bufio.NewScanner(os.Stdin)
	var urls []string
	for scanner.Scan() {
		url := scanner.Text()
		if url == "" {
			break
		}
		urls = append(urls, url)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		sig := <-sigs
		fmt.Println("SIGNAL RECEIVED:", sig)
		cancel()
	}()

	workers := createWorkers(workProcNum)
	wg := sync.WaitGroup{}
	for i := range workers{
		wg.Add(1)
		go workers[i].processUrl(&wg, ctx)
	}

	for i := range urls {
		urlsChan <- urls[i]
	}
	close(urlsChan)
	wg.Wait()
	getStatistics(workers)
}

func getStatistics(workers []Worker) {
	for i := range workers {
		fmt.Printf("Goroutine %d:%d\n", workers[i].id, workers[i].reqCount)
	}
}