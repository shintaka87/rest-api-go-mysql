CREATE DATABASE IF NOT EXISTS GOLANG_MYSQL;
 
CREATE TABLE IF NOT EXISTS GOLANG_MYSQL.users(
id INT AUTO_INCREMENT,
name VARCHAR(100),
email VARCHAR(100),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
PRIMARY KEY (id)
);

INSERT INTO users(name, email) VALUES 
("Shin", "email@example.com"),
("New User", "email2@example.com"),
("Test User", "email3@example.com");

CREATE TABLE IF NOT EXISTS GOLANG_MYSQL.posts(
id INT AUTO_INCREMENT, 
title VARCHAR(100), 
user_id INT, 
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
PRIMARY KEY (id)
); 

INSERT INTO posts(title, user_id) VALUES
("This is my first post!", 1),
("This is my Second post!", 2),
("This is my Third post", 1);


