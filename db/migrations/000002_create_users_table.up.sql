CREATE TABLE IF NOT EXISTS users(
   id uuid,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (300) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   avatar_url VARCHAR (300),
   description VARCHAR (300),
   password_remind_id VARCHAR (300),
   PRIMARY KEY(id)
);