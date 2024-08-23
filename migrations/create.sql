create table if not exists session (
    id text primary key
);

create table if not exists text (
    id text primary key,
    session text not null,
    text text,
    foreign key(session) references session(id) on delete cascade
);
