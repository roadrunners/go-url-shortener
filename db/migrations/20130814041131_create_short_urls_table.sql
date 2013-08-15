
-- +goose Up
create table short_urls (
  slug varchar(20) not null collate latin1_bin,
  url text not null,
  primary key (slug)
);

-- +goose Down
drop table short_urls;
