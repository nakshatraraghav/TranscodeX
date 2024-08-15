import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

import { S3Bucket } from "../constructs/s3";
import { SQSQueue } from "../constructs/sqs"
import { ECSCluster } from "../constructs/ecs";
import { LambdaFunction } from "../constructs/lambda"

export class LocalStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    const storage = new S3Bucket(this, "transcodex-storage-s3");

    const queue = new SQSQueue(this, "transcodex-sqs-queue");

    const cluster = new ECSCluster(this, "transcodex-worker-cluster")
    const arn = cluster.getTaskDefinition()

    const lambda = new LambdaFunction(this, "transcodex-lambda", {
      queue: queue.queue,
      taskDefinitionARN: arn
    });
  }
}
