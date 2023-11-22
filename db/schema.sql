-- Enable the use of citext (case-insensitive text)
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
	user_id serial primary key,
	email citext not null unique,
	display_name varchar(100) not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth_tokens (
	user_id int not null,
	access_token varchar(100) not null,
	refresh_token varchar(100) not null,
	expiration timestamp with time zone not null,
	FOREIGN KEY(user_id) REFERENCES users(user_id)
);


CREATE TABLE IF NOT EXISTS lists (
	list_id serial primary key,
	name varchar(100) not null,
	creator_id int not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(creator_id) REFERENCES users(user_id)
);

-- This junction table is used to model the many to many relationship between users and lists
CREATE TABLE IF NOT EXISTS user_list (
	user_id INT not null,
	list_id INT not null,
	PRIMARY KEY (user_id, list_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (list_id) REFERENCES lists(list_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tasks (
	task_id serial primary key,
	content varchar(256) not null,
	list_id int not null,
	creator_id int not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(list_id) REFERENCES lists(list_id),
	FOREIGN KEY(creator_id) REFERENCES users(user_id)
);

-- Create index for gettings tasks by list
CREATE INDEX IF NOT EXISTS idx_tasks_list_id ON tasks(list_id);

CREATE TABLE IF NOT EXISTS invitations (
	invitation_id serial primary key,
	sender_id int not null,
	receiver_id int not null,
	list_id int not null,
	FOREIGN KEY(sender_id) REFERENCES users(user_id),
	FOREIGN KEY(receiver_id) REFERENCES users(user_id),
	FOREIGN KEY(list_id) REFERENCES lists(list_id)
);

-- Create index for finding received invitations for a user
CREATE INDEX IF NOT EXISTS idx_invitations_receiver_id ON invitations(receiver_id);
