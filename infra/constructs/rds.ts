import { Construct } from "constructs";

import * as rds from "aws-cdk-lib/aws-rds";
import * as ec2 from "aws-cdk-lib/aws-ec2";

export class RDSDatabaseInstance extends Construct {
  public readonly database: rds.DatabaseInstance;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    const vpc = ec2.Vpc.fromLookup(this, "default_vpc", {
      isDefault: true,
    });

    this.database = new rds.DatabaseInstance(this, "transcodex-database-id", {
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
    });
  }
}
