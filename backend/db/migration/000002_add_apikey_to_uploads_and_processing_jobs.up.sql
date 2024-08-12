ALTER TABLE "uploads"
ADD COLUMN "apikey_id" UUID NOT NULL,
ADD CONSTRAINT "fk_uploads_api_key_id" FOREIGN KEY ("apikey_id") REFERENCES "api_keys"("id");

ALTER TABLE "processing_jobs"
ADD COLUMN "apikey_id" UUID NOT NULL,
ADD CONSTRAINT "fk_processing_jobs_api_key_id" FOREIGN KEY ("apikey_id") REFERENCES "api_keys"("id");
