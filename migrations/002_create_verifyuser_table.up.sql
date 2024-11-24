CREATE TABLE IF NOT EXISTS VerifyUsers(
   id serial PRIMARY KEY,
   username VARCHAR (50)  NOT NULL,
   verify_code VARCHAR (50)  NOT NULL,
   verify_status VARCHAR (50)  NOT NULL,
   created_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP,
   updated_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP,
);