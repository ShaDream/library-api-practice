CREATE TABLE author
(
    id   serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE publisher
(
    id   serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE reader
(
    id   serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE section
(
    id      serial PRIMARY KEY,
    shelf   int NOT NULL,
    section int NOT NULL,
    UNIQUE (shelf, section)
);

CREATE TABLE book
(
    id           serial PRIMARY KEY,
    title        varchar(255)                  NOT NULL,
    description  varchar(1023),
    author_id    int REFERENCES author (id)    NOT NULL,
    release_year smallint,
    edition      varchar(63),
    publisher_id int REFERENCES publisher (id) NOT NULL
);

CREATE TABLE physical_book
(
    id         serial PRIMARY KEY,
    book_id    int REFERENCES book (id) NOT NULL,
    section_id int REFERENCES section (id)
);

CREATE TABLE issuance
(
    id                   serial PRIMARY KEY,
    physical_book_id     int REFERENCES book (id)   NOT NULL,
    reader_id            int REFERENCES reader (id) NOT NULL,
    issue_date           date                       NOT NULL,
    complete_date        date                       NOT NULL
        CONSTRAINT complete_day_after CHECK ( complete_date >= issue_date ),
    actual_complete_date date
        CONSTRAINT actual_complete_day_after CHECK ( actual_complete_date >= issue_date )
);