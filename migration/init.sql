CREATE DATABASE wbg;
GRANT ALL PRIVILEGES ON DATABASE wbg TO almaz;
\c wbg

create table if not exists users (
                       userId uuid primary key ,
                       balance bigint not null default 0
);

create table if not exists transactions (
                              transactionId uuid primary key,
                              userId uuid references users(userId),
                              count bigint not null,
                              input bool not null ,
                              confirmed bool default false,
                              info text not null
);