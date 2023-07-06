create table customertable
(
    id       serial
        constraint customertable_pk
            primary key,
    capital  integer default 10000 not null,
    username varchar,
    password varchar
);

alter table customertable
    owner to postgres;



create table loadertable
(
    id        serial
        constraint loadertable_pk
            primary key,
    maxweight integer default 20    not null,
    drunk     boolean default false not null,
    fatigue   integer default 0     not null,
    salary    integer default 10000 not null,
    username  varchar,
    password  varchar
);

alter table loadertable
    owner to postgres;

create table tasks
(
    id     serial
        constraint tasks_pk
            primary key,
    name   varchar            not null,
    weight integer default 10 not null,
    done   integer default 0  not null
);

alter table tasks
    owner to postgres;


create table donetasks
(
    task_id   integer not null
        constraint donetasks_tasks_id_fk
            references tasks,
    loader_id integer not null
        constraint donetasks_loadertable_id_fk
            references loadertable
);

alter table donetasks
    owner to postgres;