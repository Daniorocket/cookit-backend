CREATE TABLE IF NOT EXISTS files
(
id uuid PRIMARY KEY NOT NULL,
encodedURL TEXT NOT NULL,
extension VARCHAR(10) NOT NULL
);

CREATE TABLE categories
(
id uuid PRIMARY KEY references files (id),
name VARCHAR(100) UNIQUE NOT NULL
);