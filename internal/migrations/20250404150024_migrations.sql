-- +goose Up
-- +goose StatementBegin
create table if not exists Users (
    id serial primary key,
    username varchar(255) not null unique,
    email varchar(255) not null,
    password varchar(255) not null,
    address text
);

create table if not exists Genres (
    id serial primary key,
    name varchar(255) not null unique
);

create table if not exists Authors (
    id serial primary key,
    name varchar(255) not null,
    bio text
);

create table if not exists Orders (
    id serial primary key,
    user_id int not null,
    order_date timestamp not null,
    total_amount decimal(10, 2) not null,
    foreign key (user_id) references Users(id)
);

create table if not exists Books (
    id serial primary key,
    title varchar(255) not null,
    author_id int not null,
    genre_id int not null,
    price decimal(10, 2) not null check(price > 0),
    stock_quantity int check(stock_quantity > 0),
    foreign key (author_id) references Authors(id),
    foreign key (genre_id) references Genres(id)
);

create table if not exists Orders_Books (
    order_id int not null,
    book_id int not null,
    quantity int not null check(quantity > 0),
    primary key (order_id, book_id),
    foreign key (order_id) references Orders(id) on delete cascade,
    foreign key (book_id) references Books(id) on delete cascade
);

create table if not exists Books_Authors (
    book_id int not null,
    author_id int not null,
    primary key(book_id, author_id),
    foreign key (book_id) references Books(id) on delete cascade,
    foreign key (author_id) references Authors(id) on delete cascade
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists Users
drop table if exists Genres
drop table if exists Authors
drop table if exists Orders
drop table if exists Books
drop table if exists Orders_Books
drop table if exists Books_Authors
-- +goose StatementEnd
