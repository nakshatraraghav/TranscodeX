import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
// import * as sqs from 'aws-cdk-lib/aws-sqs';

import { EC2Instance } from "../constructs/ec2";
import { S3Bucket } from "../constructs/s3";
import { RDSDatabaseInstance } from "../constructs/rds";

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const backend = new EC2Instance(this, "transcodex-backend");

    const storage = new S3Bucket(this, "transcodex-storage-s3");

    const database = new RDSDatabaseInstance(
      this,
      "transcodex-storage-database"
    );
  }
}
