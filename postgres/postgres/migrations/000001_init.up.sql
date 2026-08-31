CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name Varchar(100) unique,
        info varchar(200)
  );
