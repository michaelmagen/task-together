CREATE TABLE IF NOT EXISTS users (
	user_id VARCHAR(255) PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	verified_email BOOLEAN NOT NULL,
	name VARCHAR(255) NOT NULL,
	given_name VARCHAR(255) NOT NULL,
	family_name VARCHAR(255) NOT NULL,
	picture VARCHAR(255) NOT NULL,
	locale VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS lists (
	list_id serial primary key,
	name varchar(100) not null,
	creator_id VARCHAR(255) not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(creator_id) REFERENCES users(user_id)
);

-- This junction table is used to model the many to many relationship between users and lists
CREATE TABLE IF NOT EXISTS user_list (
	user_id VARCHAR(255) not null,
	list_id INT not null,
	PRIMARY KEY (user_id, list_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (list_id) REFERENCES lists(list_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tasks (
	task_id serial primary key,
	content varchar(256) not null,
	list_id int not null,
	creator_id VARCHAR(255) not null,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	completed boolean default false, 
	completer_id VARCHAR(255),
	FOREIGN KEY(list_id) REFERENCES lists(list_id),
	FOREIGN KEY(creator_id) REFERENCES users(user_id),
	FOREIGN KEY(completer_id) REFERENCES users(user_id)
);

-- Create index for gettings tasks by list
CREATE INDEX IF NOT EXISTS idx_tasks_list_id ON tasks(list_id);

CREATE TABLE IF NOT EXISTS invitations (
	invitation_id serial primary key,
	sender_id  VARCHAR(255) not null,
	receiver_id VARCHAR(255) not null,
	list_id int not null,
	FOREIGN KEY(sender_id) REFERENCES users(user_id),
	FOREIGN KEY(receiver_id) REFERENCES users(user_id),
	FOREIGN KEY(list_id) REFERENCES lists(list_id)
);

-- Create index for finding received invitations for a user
CREATE INDEX IF NOT EXISTS idx_invitations_receiver_id ON invitations(receiver_id);
