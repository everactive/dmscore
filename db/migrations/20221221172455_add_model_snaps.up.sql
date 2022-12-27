CREATE TABLE "device_models" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"name" text,PRIMARY KEY ("id"));

CREATE INDEX "idx_device_models_deleted_at" ON "device_models" ("deleted_at");

CREATE TABLE "device_model_required_snaps" ("id" bigserial,"created_at" timestamptz,"updated_at" timestamptz,"deleted_at" timestamptz,"device_model_id" bigint,"name" text,PRIMARY KEY ("id"),CONSTRAINT "fk_device_models_device_model_required_snaps" FOREIGN KEY ("device_model_id") REFERENCES "device_models"("id"));

CREATE INDEX "idx_device_model_required_snaps_deleted_at" ON "device_model_required_snaps" ("deleted_at");

CREATE UNIQUE INDEX "idx_device_model_required_snaps_name" ON "device_model_required_snaps" ("name")