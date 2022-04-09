-- Adminer 4.8.1 PostgreSQL 14.2 (Debian 14.2-1.pgdg110+1) dump

DROP TABLE IF EXISTS "nasabah";
DROP SEQUENCE IF EXISTS nasabah_id_seq;
CREATE SEQUENCE nasabah_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."nasabah" (
    "id" integer DEFAULT nextval('nasabah_id_seq') NOT NULL,
    "username" character varying(250) NOT NULL,
    "password" character varying(250) NOT NULL,
    "token" text NOT NULL,
    "tabungan" integer NOT NULL,
    CONSTRAINT "nasabah_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "nasabah_username_key" UNIQUE ("username")
) WITH (oids = false);

INSERT INTO "nasabah" ("id", "username", "password", "token", "tabungan") VALUES
(1,	'hero',	'hero',	'token',	10000),
(2,	'van helsing',	'van helsing',	'token',	35000);

-- 2022-04-09 06:17:19.945812+00
