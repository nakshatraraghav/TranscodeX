#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { InfraStack } from "../lib/infra-stack";

import { config } from "dotenv";
import { env } from "../config/zenv"
config({
  path: "../.env"
});

const app = new cdk.App();
new InfraStack(app, "InfraStack", {
  env: {
    account: env.AWS_ACCOUNT_ID,
    region: env.AWS_REGION,
  },
});
