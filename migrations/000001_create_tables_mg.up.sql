DROP TABLE IF EXISTS users CASCADE;


CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(64) NOT NULL,
    phone_number varchar(20) NOT NULL UNIQUE
);

DROP TABLE IF EXISTS tasks CASCADE;

CREATE TABLE tasks
(
    task_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title varchar(64) NOT NULL,
    description varchar(255),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    isDone BOOLEAN NOT NULL,
    deletedAt TIMESTAMP DEFAULT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


