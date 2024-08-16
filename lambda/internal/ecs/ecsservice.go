package ecsservice

import (
	"fmt"

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

	// Check if the required attributes are present
	mediaType, ok := record.MessageAttributes["media_type"]
	if !ok {
		return fmt.Errorf("missing required attribute: media_type")
	}

	uploadID, ok := record.MessageAttributes["upload_id"]
	if !ok {
		return fmt.Errorf("missing required attribute: upload_id")
	}

	objectKey, ok := record.MessageAttributes["object_key"]
	if !ok {
		return fmt.Errorf("missing required attribute: object_key")
	}

	transformations, ok := record.MessageAttributes["transformations"]
	if !ok {
		return fmt.Errorf("missing required attribute: transformations")
	}

	fmt.Println("PRINTTTTTT", env)

	// Run the ECS task
	_, err := e.client.RunTask(&ecs.RunTaskInput{
		Cluster:        aws.String(env.ECS_CLUSTER_NAME),
		TaskDefinition: aws.String(env.ECS_TASK_DEFINITION),
		LaunchType:     aws.String("FARGATE"),
		NetworkConfiguration: &ecs.NetworkConfiguration{
			AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
				Subnets: aws.StringSlice(env.SUBNET_IDS),
				SecurityGroups: []*string{
					&env.SECURITY_GROUP_ID,
				},
			},
		},
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				{
					Name: aws.String("transcodex-worker"),
					Environment: []*ecs.KeyValuePair{
						{
							Name:  aws.String("MEDIA_TYPE"),
							Value: mediaType.StringValue,
						},
						{
							Name:  aws.String("BUCKET_NAME"),
							Value: aws.String(env.BUCKET_NAME),
						},
						{
							Name:  aws.String("OBJECT_KEY"),
							Value: objectKey.StringValue,
						},
						{
							Name:  aws.String("TRANSFORMATIONS"),
							Value: transformations.StringValue,
						},
						{
							Name:  aws.String("CONNECTION_STRING"),
							Value: aws.String(CONNECTION_STRING),
						},
						{
							Name:  aws.String("UPLOAD_ID"),
							Value: uploadID.StringValue,
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
