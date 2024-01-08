CREATE TABLE vessels
(
    imo int NOT NULL,
    name varchar(255) NOT NULL,
    flag varchar(2) NOT NULL,
    year_built int NOT NULL,
    owner varchar(255) NOT NULL,

    PRIMARY KEY (imo)
)