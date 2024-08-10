import { Construct } from "constructs";

import * as s3 from "aws-cdk-lib/aws-s3";

export class S3Bucket extends Construct {
  public readonly bucket: s3.Bucket;

  constructor(scope: Construct, id: string) {
    super(scope, id);

    this.bucket = new s3.Bucket(this, "transcodex-s3-bucket-id", {
      bucketName: "storage.bucket.transcodex",
      objectOwnership: s3.ObjectOwnership.BUCKET_OWNER_ENFORCED,
      blockPublicAccess: s3.BlockPublicAccess.BLOCK_ALL,
      versioned: true,
      encryption: s3.BucketEncryption.S3_MANAGED,
    });
  }
}
