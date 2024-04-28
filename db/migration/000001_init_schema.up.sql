CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SET TIMEZONE='Asia/Ho_Chi_Minh';

CREATE TABLE "accounts" (
                            "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
                            "owner" varchar NOT NULL,
                            "balance" bigint NOT NULL,
                            "currency" varchar NOT NULL,
                            "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
                           "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
                           "account_id" uuid NOT NULL,
                           "amount" bigint NOT NULL,
                           "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
                             "id" uuid PRIMARY KEY DEFAULT (uuid_generate_v4()),
                             "from_account_id" uuid NOT NULL,
                             "to_account_id" uuid NOT NULL,
                             "amount" bigint NOT NULL,
                             "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
