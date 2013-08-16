
-- +goose Up
create table short_urls (
  id int not null auto_increment,
  url text not null,
  primary key (id)
);

-- +goose Down
drop table short_urls;
