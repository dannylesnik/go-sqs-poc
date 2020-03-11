package queueing

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

//StartPolling -
func (cvs *SQS) StartPolling(chn chan<- *sqs.Message) {

	for {
		receiveparams := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(os.Getenv("QUEUE_URL")),
			MaxNumberOfMessages: aws.Int64(1),
			VisibilityTimeout:   aws.Int64(360),
			WaitTimeSeconds:     aws.Int64(20),
		}

		receiveresp, err := cvs.ReceiveMessage(receiveparams)
		if err != nil {
			log.Println(err)
		}

		for _, message := range receiveresp.Messages {
			fmt.Printf("[Receive message] sent to channel \n%v \n\n", *message.MessageId)
			chn <- message
		}
	}
}

//DeleteMSG -
func (cvs *SQS) DeleteMSG(msg *sqs.Message) error {

	deleteparams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(os.Getenv("QUEUE_URL")), // Required
		ReceiptHandle: msg.ReceiptHandle,                  // Required

	}
	_, err := cvs.DeleteMessage(deleteparams) // No response returned when successed.
	return err
}
