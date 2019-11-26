BEGIN;
--
-- Create entities User
--
CREATE TABLE IF NOT EXISTS "ip_list"
(
    "id"           serial                   NOT NULL PRIMARY KEY,
    "ip"           CIDR UNIQUE,
    "kind"         varchar(10)              NOT NULL CHECK ("kind" IN ('black', 'white')),
    "date_created" timestamp with time zone NOT NULL
);
CREATE INDEX IF NOT EXISTS list_ip ON "ip_list" ("ip");
COMMIT;
