import { Construct } from "constructs";

import * as sqs from "aws-cdk-lib/aws-sqs";
import { Duration } from "aws-cdk-lib";

export class SQSQueue extends Construct {
  public readonly queue: sqs.Queue;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    this.queue = new sqs.Queue(this, "transcodex-queue-id", {
      queueName: "transcodex-sqs-queue",
      fifo: false,
      visibilityTimeout: Duration.minutes(1),
      encryption: sqs.QueueEncryption.SQS_MANAGED,
    });
  }
}
