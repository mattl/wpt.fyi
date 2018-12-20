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

func startMon(maxReqs uint, timeout time.Duration) chan int {
	log.Infof("Monitor started")

	req := make(chan int)
	go func() {
		var timer chan bool
		var reqs uint
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
				reqs = 0
			case <-req:
				reqs++
				if reqs == maxReqs {
					log.Infof("Monitor: Max reqs")
					time.Sleep(time.Second)
					reqCounter.reset()
					time.Sleep(time.Second)
					reqs = 0
				}
			}
		}
	}()

	return req
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
			log.Infof("Starting req %d", j)
			req <- j
			log.Infof("Running req %d", j)
			reqCounter.inc()
		}(i)
	}
	wg.Wait()

	time.Sleep(time.Second * 20)
}
