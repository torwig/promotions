CREATE TABLE promotions (
    id CHAR(36),
    price FLOAT NOT NULL,
    expiration_date DATETIME NOT NULL,
    PRIMARY KEY(id)
);
