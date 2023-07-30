CREATE TABLE _admins (
  ID INTEGER PRIMARY KEY,
  email varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL
);

-- type 
-- 0 => base
-- 1 => auth
CREATE TABLE _collections (
  ID INTEGER PRIMARY KEY,
  name varchar(255) NOT NULL,
  type integer NOT NULL
);

CREATE TABLE users (
  ID INTEGER PRIMARY KEY,
  username varchar(255) NOT NULL,
  email varchar(255) UNIQUE NOT NULL,
  password varchar(255) NOT NULL
);
