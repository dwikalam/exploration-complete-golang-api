CREATE TABLE IF NOT EXISTS "user" (
    "id" SERIAL,
    "fullname" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) UNIQUE NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    PRIMARY KEY ("id")
);

CREATE OR REPLACE TRIGGER update_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
        EXECUTE FUNCTION update_updated_at_column();