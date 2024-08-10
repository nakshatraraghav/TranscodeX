import { Construct } from "constructs";

import * as ecs from "aws-cdk-lib/aws-ecs";

export class ECSCluster extends Construct {
  public readonly cluster: ecs.Cluster;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    this.cluster = new ecs.Cluster(this, "transcodex-cluster-id", {})

    const audio = new ecs.TaskDefinition(this, "transcodex-audio-worker-id", {
      compatibility: ecs.Compatibility.FARGATE
    })

    const video = new ecs.TaskDefinition(this, "transcodex-video-worker-id", {
      compatibility: ecs.Compatibility.FARGATE
    })

    const image = new ecs.TaskDefinition(this, "transcodex-image-worker-id", {
      compatibility: ecs.Compatibility.FARGATE
    })

    
  }
}