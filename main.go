package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-sqs-poc/queueing"
)

//Env -
type Env struct {
	SQS         *queueing.SQS
	ChnMessages chan *sqs.Message
}

func (env *Env) handleMessage(counter int) {

	fmt.Printf("Handler with id [%d] started.\n\n", counter)
	for message := range env.ChnMessages {
		fmt.Printf("[Receive message] by thread [%d] received from channel \n%v \n\n", counter, *message.MessageId)
	    time.Sleep(5 * time.Minute)

		err := env.SQS.DeleteMSG(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("[Delete message] \nMessage ID: %s has beed deleted by thread [%d].\n\n", *message.MessageId, counter)
		}
	}
	fmt.Printf("Handler with id [%d] ended\n\n", counter)

}

func main() {

	svc := queueing.InitSQS()
	chnMessages := make(chan *sqs.Message, 3)
	go svc.StartPolling(chnMessages)

	env := &Env{svc, chnMessages}

	for i := 1; i < 4; i++ {
		go env.handleMessage(i)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go env.waitForShutdown(&wg)
	wg.Wait()
}

func (env *Env) waitForShutdown(wg *sync.WaitGroup) {

	defer wg.Done()
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	close(env.ChnMessages)
	// Create a deadline to wait for.
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Println("Shutting down")

}
