CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" VARCHAR NOT NULL,
  "balance" BIGINT NOT NULL,
  "currency" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" BIGSERIAL PRIMARY KEY,
  "to_account" BIGINT NOT NULL,
  "from_account" BIGINT NOT NULL,
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transactions" ("to_account");

CREATE INDEX ON "transactions" ("from_account");

CREATE INDEX ON "transactions" ("to_account", "from_account");

COMMENT ON COLUMN "entries"."amount" IS 'Can be negative';

COMMENT ON COLUMN "transactions"."amount" IS 'Can not be neagative';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_account") REFERENCES "accounts" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_account") REFERENCES "accounts" ("id");
