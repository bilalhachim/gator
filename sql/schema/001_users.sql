-- +goose Up
CREATE TABLE users (
 id  UUID PRIMARY KEY NOT NULL,
 created_at TIMESTAMP NOT NULL,
 updated_at TIMESTAMP NOT NULL,
 name varchar(255) NOT NULL,
 UNIQUE(name) 
);
CREATE TABLE feeds( 
feed_id UUID PRIMARY KEY NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL, 
name varchar(255) NOT NULL, 
url varchar(255) NOT NULL, 
UNIQUE(url),
reference_id UUID NOT NULL,
last_fetched_at TIMESTAMP, 
FOREIGN KEY (REFERENCE_ID) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE feed_follows (
id  UUID PRIMARY KEY NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
user_id UUID NOT NULL,
feed_id UUID NOT NULL,
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE,
UNIQUE(user_id,feed_id)
 );
CREATE TABLE posts(
id  UUID PRIMARY KEY NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
title varchar(255) NOT NULL, 
url varchar(255) NOT NULL, 
UNIQUE(url),
description varchar(255) NOT NULL, 
published_at TIMESTAMP NOT NULL,
feed_id UUID NOT NULL,
FOREIGN KEY (feed_id) REFERENCES feeds(feed_id) ON DELETE CASCADE
 );


-- +goose Down
DROP TABLE feed_follows;
DROP TABLE posts;
DROP TABLE feeds;
DROP TABLE users;
