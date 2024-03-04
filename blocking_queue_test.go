package bqueue

import (
	"fmt"
	"testing"
	"time"
)

var bq = NewBlockingQueue(10)

func TestBlockingQueue_Timeout(t *testing.T) {

	go func() {
		for {
			// consumer, 800ms timeout
			_, ok := bq.Poll(800 * time.Millisecond)
			if ok {
				fmt.Printf("poll success -> %v\n", time.Now().UnixMilli())
			} else {
				fmt.Printf("poll timeout -> %v\n", time.Now().UnixMilli())
			}
		}
	}()

	// producer, send rate 2s
	tk := time.NewTicker(2 * time.Second)
	go func() {
		for {
			<-tk.C
			// producer, 1000ms timeout
			offered := bq.Offer(1, time.Second)
			if !offered {
				fmt.Println("enqueue timeout")
			} else {
				fmt.Println("enqueue success")
			}
		}
	}()

	<-time.After(time.Second * 30)
	tk.Stop()
}
