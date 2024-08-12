
ALTER TABLE "uploads"
DROP CONSTRAINT "fk_uploads_api_key_id",
DROP COLUMN "apikey_id";

ALTER TABLE "processing_jobs"
DROP CONSTRAINT "fk_processing_jobs_api_key_id",
DROP COLUMN "apikey_id";
