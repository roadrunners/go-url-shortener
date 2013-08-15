
-- +goose Up
create table short_urls (
  slug varchar(20) not null,
  url text not null,
  primary key (slug)
);

-- +goose Down
drop table short_urls;
