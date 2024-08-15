package ecsservice

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/nakshatraraghav/transcodex/lambda/config"
)

type ECSService struct {
	client *ecs.ECS
}

func NewECSService(session *session.Session) *ECSService {
	client := ecs.New(session)

	return &ECSService{
		client: client,
	}
}

func (e *ECSService) RunTask(record events.SQSMessage, CONNECTION_STRING string) error {

	env := config.Getenv()

	mediaType := record.Attributes["media_type"]
	uploadID := record.Attributes["upload_id"]
	objectKey := record.Attributes["object_key"]
	transformations := record.Attributes["transformations"]

	_, err := e.client.RunTask(&ecs.RunTaskInput{
		Cluster:        aws.String(env.ECS_CLUSTER_NAME),
		TaskDefinition: aws.String(env.ECS_TASK_DEFINITION),
		LaunchType:     aws.String("FARGATE"),
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				&ecs.ContainerOverride{
					Name: aws.String("transcodex-worker"),
					Environment: []*ecs.KeyValuePair{
						&ecs.KeyValuePair{
							Name:  aws.String("MEDIA_TYPE"),
							Value: aws.String(mediaType),
						},
						&ecs.KeyValuePair{
							Name:  aws.String("BUCKET_NAME"),
							Value: aws.String(env.BUCKET_NAME),
						},
						&ecs.KeyValuePair{
							Name:  aws.String("OBJECT_KEY"),
							Value: aws.String(objectKey),
						},
						&ecs.KeyValuePair{
							Name:  aws.String("TRANSFORMATIONS"),
							Value: aws.String(transformations),
						},
						&ecs.KeyValuePair{
							Name:  aws.String("CONNECTION_STRING"),
							Value: aws.String(CONNECTION_STRING),
						},
						&ecs.KeyValuePair{
							Name:  aws.String("UPLOAD_ID"),
							Value: aws.String(uploadID),
						},
					},
				},
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}
