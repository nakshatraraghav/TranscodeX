import { Construct } from "constructs";

import * as lambda from "aws-cdk-lib/aws-lambda";
import * as iam from "aws-cdk-lib/aws-iam"
import * as sqs from "aws-cdk-lib/aws-sqs"
import * as lambdaEventSources from "aws-cdk-lib/aws-lambda-event-sources";

import path = require("path");

import { env } from "../config/zenv";

interface LambdaFunctionProps {
  queue: sqs.IQueue
  taskDefinitionARN: string
}

export class LambdaFunction extends Construct {
  public readonly func: lambda.Function;
  
  constructor(scope: Construct, id: string, props: LambdaFunctionProps) {
    super(scope, id)

    const lambdaRole = new iam.Role(this, "lambda-execution-role-id", {
      assumedBy: new iam.ServicePrincipal("lambda.amazonaws.com")
    })

    lambdaRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonRDSFullAccess")
    );
    lambdaRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonSQSFullAccess")
    );
    lambdaRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonECS_FullAccess")
    );

    const p = path.join(__dirname, "..", "..", "lambda", "bin", "lambda.zip")

    this.func = new lambda.Function(this, "transcodex-lambda-id", {
      runtime: lambda.Runtime.PROVIDED_AL2,
      code: lambda.Code.fromAsset(p),
      handler: "main",
      environment: {
        REGION_STRING: env.AWS_REGION,
        BUCKET_NAME: env.BUCKET_NAME,
        ECS_CLUSTER_NAME: env.ECS_CLUSTER_NAME,
        ECS_TASK_DEFINITION: props.taskDefinitionARN,
        RDS_DATABASE_USERNAME: env.RDS_DATABASE_USERNAME,
        RDS_DATABASE_PASSWORD: env.RDS_DATABASE_PASSWORD,
        DATABASE_INSTANCE_IDENTIFIER: env.DATABASE_INSTANCE_IDENTIFIER
      },
      role: lambdaRole
    })

    this.func.addEventSource(new lambdaEventSources.SqsEventSource(props.queue, {
      batchSize: 2
    }))

  }
}