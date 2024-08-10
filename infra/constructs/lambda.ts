import { Construct } from "constructs";

import * as lambda from "aws-cdk-lib/aws-lambda";

export class LambdaFunction extends Construct {
  public readonly func: lambda.Function;
  
  constructor(scope: Construct, id: string) {
    super(scope, id)
  }
}