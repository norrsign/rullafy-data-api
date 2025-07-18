CREATE TABLE users (
  id       text PRIMARY KEY,
  job      text NOT NULL
);

CREATE TABLE companies (
  id       text PRIMARY KEY,
  name     text NOT NULL,
  address  text,
  phone    text
);


 

CREATE TABLE products (
  id       text PRIMARY KEY,
  name     text NOT NULL,
  type     text NOT NULL,
  description text NOT NULL
);



