package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Saloon capacity, number of clients, number of barbers
var saloonCapacity int = 1
var clientCount int = 10
var barberCount int = 1

// Barbers semaphore, saloon channel used both for synchronisation and queueing clients,
// waiting variable keeps track of clients in the saloon, mutex for critical section
var barbers = make(chan int)
var saloon = make(chan int, saloonCapacity)
var waiting int = 0
var mutex sync.Mutex

// WaitGroups to avoid deadlock
var cWg sync.WaitGroup
var bWg sync.WaitGroup

func sleep() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
}

func barber(id int) {
	defer bWg.Done()

	for client := range saloon {
		// One customer is served
		barbers <- 1
		fmt.Printf("Barber%d: cuts Client%d\n", id, client)
		sleep()
	}
}

func client(id int) {
	defer cWg.Done()

	mutex.Lock()
	if waiting < saloonCapacity {
		// Client enters saloon
		waiting++
		saloon <- id
		fmt.Printf("\t\t\tClient%d: enters saloon\n", id)
		mutex.Unlock()

		// Waits to be served
		<-barbers

		// Client leaves saloon
		mutex.Lock()
		waiting--
		mutex.Unlock()

	} else {
		// Saloon is full, client leaves
		fmt.Printf("\t\t\tClient%d: saloon is full, leaves\n", id)
		mutex.Unlock()
	}
}

func main() {

	// Launch barber goroutines
	for i := 0; i < barberCount; i++ {
		bWg.Add(1)
		go barber(i)
	}

	// Launch client goroutines
	for i := 0; i < clientCount; i++ {
		cWg.Add(1)
		go client(i)
		sleep()
	}

	// Wait for client gorotines to finish, close saloon channel, wait for barbers
	cWg.Wait()
	close(saloon)
	bWg.Wait()
}
