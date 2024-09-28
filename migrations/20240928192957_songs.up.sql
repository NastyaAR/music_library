create table if not exists songs (
    song_group text,
    name text,
    release_date date,
    text text,
    link text,
    primary key (song_group, name)
);