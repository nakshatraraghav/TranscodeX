import { Construct } from "constructs";

import * as rds from "aws-cdk-lib/aws-rds";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import { SecretValue } from "aws-cdk-lib";

import { env } from "../config/zenv";

export class RDSDatabaseInstance extends Construct {
  public readonly database: rds.DatabaseInstance;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    const vpc = ec2.Vpc.fromLookup(this, "default_vpc", {
      isDefault: true,
    });

    const password = SecretValue.unsafePlainText(env.RDS_DATABASE_PASSWORD as string)

    this.database = new rds.DatabaseInstance(this, "transcodex-database-id", {
      instanceIdentifier: env.DATABASE_INSTANCE_IDENTIFIER,
      databaseName: "transcodex-database",
      engine: rds.DatabaseInstanceEngine.postgres({
        version: rds.PostgresEngineVersion.VER_16,
      }),
      vpc: vpc,
      instanceType: ec2.InstanceType.of(
        ec2.InstanceClass.T3,
        ec2.InstanceSize.MICRO
      ),
      multiAz: false,
      publiclyAccessible: false,
      credentials: rds.Credentials.fromPassword(env.RDS_DATABASE_USERNAME as string, password)
    });
  }
}

// RDS_DATABASE_PASSWORD
// RDS_DATABASE_USERNAME