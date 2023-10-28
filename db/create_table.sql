CREATE TABLE images.images(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    title VARCHAR(50),
    post_user VARCHAR(50),
    image_path VARCHAR(50),
    category VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
