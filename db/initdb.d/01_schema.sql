create table users (
  id bigint,
  primary key id (id)
);

create table posts (
  id bigint,
  author_id bigint not null,
  content text not null,
  primary key id (id),
  constraint fk_author_id foreign key (author_id) references users (id)
);
