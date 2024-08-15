import { z } from "zod";
import { config } from "dotenv"

const schema = z.object({
  AWS_ACCOUNT_ID: z.string(), //
  AWS_REGION: z.string(), //
  TRANSCODEX_WORKER_IMAGE_URI: z.string(), //
  DATABASE_INSTANCE_IDENTIFIER: z.string(), //
  RDS_DATABASE_USERNAME: z.string(), //
  RDS_DATABASE_PASSWORD: z.string(), //
  BUCKET_NAME: z.string(), //
  ECS_CLUSTER_NAME: z.string(),
});

function loadenv() {
  config();

  const env = schema.parse(process.env);

  return env;
}

export const env = loadenv();

