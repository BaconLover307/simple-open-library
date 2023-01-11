CREATE TABLE book
(
    book_id varchar(20) PRIMARY KEY,
    title varchar(200),
    edition integer
);

CREATE TABLE pickup(
    pickup_id integer PRIMARY KEY AUTO_INCREMENT,
    book_id varchar(20),
    schedule datetime,
    FOREIGN KEY (book_id) REFERENCES book(book_id)
);

CREATE TABLE author(
   author_id varchar(20) PRIMARY KEY,
   name varchar(100)
);

CREATE TABLE authored(
   author_id varchar(20),
   book_id varchar(20),
   PRIMARY KEY (author_id, book_id),
   FOREIGN KEY (author_id) REFERENCES author(author_id),
   FOREIGN KEY (book_id) REFERENCES book(book_id)
);
