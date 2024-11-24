CREATE TABLE IF NOT EXISTS Users(
   id serial PRIMARY KEY,
   name VARCHAR (50)  NOT NULL,
   username VARCHAR (50) UNIQUE NOT NULL,
   email VARCHAR (50) UNIQUE NOT NULL,
   phone VARCHAR (50) UNIQUE NOT NULL,
   password text NOT NULL,
   token text NOT NULL,
   created_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP,
   updated_at timestamptz (6) DEFAULT CURRENT_TIMESTAMP
);