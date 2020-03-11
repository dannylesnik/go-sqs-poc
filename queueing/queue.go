package queueing

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const (
	//Region - SQS Queue Path
	Region = "us-east-1"
	//CredPath - SQS Queue Path
	CredPath = "/Users/dannylesnik/.aws/credentials"
	//CredProfile - SQS Queue Path
	CredProfile = "trtms"
)

//SQS -sqs warapper
type SQS struct {
	*sqs.SQS
}

//InitSQS -
func InitSQS() *SQS {
	sess := session.New(&aws.Config{
		Region:      aws.String(Region),
		Credentials: credentials.NewSharedCredentials(CredPath, CredProfile),
		MaxRetries:  aws.Int(5),
	})

	svc := sqs.New(sess)
	return &SQS{svc}
}
