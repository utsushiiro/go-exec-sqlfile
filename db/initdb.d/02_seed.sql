insert into users (id) values (1);
insert into users (id) values (2);

insert into posts (id, author_id, content) values (1, 1, 'Hello, world!');
-- insert into posts (id, author_id, content) values (2, 1, 'Comment Out');
insert into posts (id, author_id, content) values (2, 1, '(;_;)');
insert into posts (id, author_id, content) values (3, 2, 'Hello, world!');
/*
 * insert into posts (id, author_id, content) values (3, 2, 'Comment Out');
 */
insert into posts (id, author_id, content) values (4, 2, 'insert into posts (id, author_id, content) values (4, 2, ''Hello, world!'');');
