CREATE TABLE IF NOT EXISTS VerifyUsers(
   id serial PRIMARY KEY,
   user_id VARCHAR (50)  NOT NULL,
   verify_code VARCHAR (50)  NOT NULL,
   verify_status VARCHAR (50)  NOT NULL,
   created_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP,
   updated_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP
   );