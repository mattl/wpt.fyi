package main

import (
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type counter struct {
	count uint
	limit uint

	m *sync.Mutex
}

func (c *counter) inc() {
	c.m.Lock()
	defer c.m.Unlock()

	c.count++
	if c.count > c.limit {
		panic("Count exceeded limit")
	}
}

func (c *counter) reset() {
	c.m.Lock()
	defer c.m.Unlock()

	c.count = 0
}

const maxReqs = 100

var reqCounter = &counter{
	count: 0,
	limit: maxReqs,
	m:     &sync.Mutex{},
}

func manageReqs(mon chan bool, reqBlock chan bool, maxReqs uint) chan int {
	log.Infof("Request manager started")

	req := make(chan int)
	go func() {
		count := uint(0)
		for {
			if count < maxReqs {
				select {
				case n := <-req:
					log.Infof("Request manager: Accepted request %d", n)
					count++
				case <-mon:
					log.Infof("Request manager: Monitor run complete")
					count = uint(0)
				}
			} else {
				log.Infof("Request manager: Blocking on request backlog")
				reqBlock <- true
				<-mon
				count = uint(0)
			}
		}
	}()

	return req
}

func runMon(mon chan bool, reqBlock chan bool, timeout time.Duration) {
	var timer chan bool
	for {
		timer = make(chan bool, 1)
		go func() {
			time.Sleep(timeout)
			timer <- true
		}()
		select {
		case <-timer:
			log.Infof("Monitor: Timed out")
			time.Sleep(time.Second)
			reqCounter.reset()
			time.Sleep(time.Second)
		case <-reqBlock:
			log.Infof("Monitor: Requests blocked on monitor")
			time.Sleep(time.Second)
			reqCounter.reset()
			time.Sleep(time.Second)
		}
		log.Infof("Monitor: Monitor run complete")
		mon <- true
	}
}

func startMon(maxReqs uint, timeout time.Duration) chan int {
	log.Infof("Monitor started")

	mon := make(chan bool)
	reqBlock := make(chan bool)
	go runMon(mon, reqBlock, timeout)

	return manageReqs(mon, reqBlock, maxReqs)
}

func main() {
	req := startMon(maxReqs, time.Second*2)

	time.Sleep(time.Second * 5)

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			r := rand.Intn(20000)
			time.Sleep(time.Millisecond * time.Duration(r))
			req <- j
			reqCounter.inc()
		}(i)
	}
	wg.Wait()

	time.Sleep(time.Second * 20)
}
