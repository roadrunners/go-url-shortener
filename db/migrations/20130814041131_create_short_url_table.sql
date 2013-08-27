
-- +goose Up
create table short_url (
  id int not null auto_increment,
  url text not null,
  primary key (id)
);

-- +goose Down
drop table short_url;
