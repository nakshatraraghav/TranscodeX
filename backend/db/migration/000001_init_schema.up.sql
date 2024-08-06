CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" VARCHAR NOT NULL,
  "username" VARCHAR UNIQUE NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "valid" BOOLEAN DEFAULT true,
  "user_agent" VARCHAR NOT NULL,
  "ip" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "api_keys" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "enabled" BOOLEAN DEFAULT true,
  "key" VARCHAR NOT NULL,
  "user_id" UUID NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "uploads" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "file_name" VARCHAR NOT NULL,
  "file_type" VARCHAR NOT NULL,
  "s3_url" VARCHAR NOT NULL,
  "status" VARCHAR NOT NULL DEFAULT 'PENDING',
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "processing_jobs" (
  "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user_id" UUID NOT NULL,
  "upload_id" UUID NOT NULL,
  "job_type" VARCHAR NOT NULL,
  "status" VARCHAR NOT NULL DEFAULT 'NOT_STARTED',
  "result_url" VARCHAR,
  "created_at" TIMESTAMPTZ DEFAULT now(),
  "updated_at" TIMESTAMPTZ DEFAULT now(),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  FOREIGN KEY ("upload_id") REFERENCES "uploads" ("id") ON DELETE CASCADE
);
