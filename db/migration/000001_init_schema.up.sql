CREATE TABLE "users" (
    "id"            BIGSERIAL PRIMARY KEY,
    "name"          VARCHAR(40) NOT NULL,
    "email"         VARCHAR(40) NOT NULL,
    "password"      VARCHAR(40) NOT NULL,
    "activated"     BOOLEAN NOT NULL,
    "created_at"    TIMESTAMP NOT NULL DEFAULT(now())
);