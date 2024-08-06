-- Drop Indexes
DROP INDEX IF EXISTS "idx_upload_id";

DROP INDEX IF EXISTS "idx_user_id";

DROP INDEX IF EXISTS "idx_user_id_ip";

DROP INDEX IF EXISTS "idx_username_email";

-- Drop Tables
DROP TABLE IF EXISTS "processing_jobs";

DROP TABLE IF EXISTS "uploads";

DROP TABLE IF EXISTS "api_keys";

DROP TABLE IF EXISTS "sessions";

DROP TABLE IF EXISTS "users";
