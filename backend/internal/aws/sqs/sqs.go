package sqsservice

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/nakshatraraghav/transcodex/backend/internal/config"
)

type SQSService struct {
	queue *sqs.SQS
}

func NewSQSService() (*SQSService, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(config.GetEnv().AWS_REGION),
	})
	if err != nil {
		return nil, err
	}

	q := sqs.New(session)

	return &SQSService{
		queue: q,
	}, nil
}

func (s *SQSService) Enqueue(mediaType, key, transformations string) error {

	_, err := s.queue.SendMessage(
		&sqs.SendMessageInput{
			QueueUrl:     aws.String(config.GetEnv().SQS_QUEUE_URL),
			MessageBody:  aws.String(fmt.Sprintf("processing %v", mediaType)),
			DelaySeconds: aws.Int64(0),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"media_type": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(mediaType),
				},
				"object_type": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(key),
				},
				"transformations": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String(transformations),
				},
			},
		},
	)

	if err != nil {
		return err
	}

	return nil
}
