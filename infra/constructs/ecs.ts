import { Construct } from "constructs";

import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam"

import { env } from "../config/zenv";

export class ECSCluster extends Construct {
  public readonly cluster: ecs.Cluster;
  public readonly task: ecs.TaskDefinition

  constructor(scope: Construct, id: string) {
    super(scope, id);

    const vpc = ec2.Vpc.fromLookup(this, "DefaultVPC", {
      vpcId: env.VPC_ID
    });

    this.cluster = new ecs.Cluster(this, "transcodex-cluster-id", {
      clusterName: env.ECS_CLUSTER_NAME,
      vpc: vpc,
      enableFargateCapacityProviders: true
    });

    const taskRole = new iam.Role(this, "transcodex-worker-task-role-id", {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com")
    });

    taskRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonRDSFullAccess")
    );

    taskRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("AmazonSQSFullAccess")
    );

    const taskExecutionRole = new iam.Role(this, "transcodex-worker-task-execution-role-id", {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
    })

    taskExecutionRole.addManagedPolicy(
      iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AmazonECSTaskExecutionRolePolicy")
    );

    this.task = new ecs.FargateTaskDefinition(this, "TranscodexTaskDefinition", {
      family: "transcodex-worker-task",
      cpu: 4096,
      memoryLimitMiB: 8192,
      taskRole: taskRole,
      executionRole: taskExecutionRole
    });

    const container = this.task.addContainer("transcodex-worker-container-id", {
      containerName: "transcodex-worker",
      image: ecs.ContainerImage.fromRegistry(env.TRANSCODEX_WORKER_IMAGE_URI as string),
      essential: true,
      logging: ecs.LogDriver.awsLogs({
        streamPrefix: "transcodex-worker-logs"
      })
    });

  }

  public getTaskDefinition(): string {
    return this.task.taskDefinitionArn
  }
}