package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Worker struct {
	id int
	reqCount int
	reqClient http.Client
}

func (w *Worker) createClient()  {
	w.reqClient = http.Client{
		Timeout: 10*time.Second,
	}
}

func (w *Worker) processUrl(wg *sync.WaitGroup, ctx context.Context)  {
	select {
	case <-ctx.Done():
		break
	case url := <-urlsChan:
		startTime := time.Now()
		res, err := w.reqClient.Get(url)
		w.reqCount += 1
		endTime := time.Now()
		if err != nil {
			fmt.Printf("Can't make request to %s error appeared: %s\n", url, err)
			break
		}
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("%s;%d;%d;%dms\n",url, res.StatusCode, len(body), endTime.Sub(startTime).Milliseconds())
	}
	wg.Done()
}

func createWorkers(numWorkers int) []Worker {
	var workers []Worker
	for i := 1; i <= numWorkers; i++ {
		w := Worker{
			id: i,
		}
		w.createClient()
		workers = append(workers, w)
	}
	return workers
}
