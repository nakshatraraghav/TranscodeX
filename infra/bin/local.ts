#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { LocalStack } from "../lib/local-stack";

import { config } from "dotenv";
import { env } from "../config/zenv"
config({
  path: "../.env"
});

const app = new cdk.App();
new LocalStack(app, "LocalStack", {
  env: {
    account: env.AWS_ACCOUNT_ID,
    region: env.AWS_REGION,
  },
});
