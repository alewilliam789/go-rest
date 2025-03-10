CREATE TABLE user (
  user_id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  user_name VARCHAR(100) NOT NULL UNIQUE,
  password VARBINARY(100) NOT NULL,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  dob VARCHAR(20),
  city VARCHAR(50), 
  state VARCHAR(2)
)
