create table customertable
(
    ID       int auto_increment
        primary key,
    capital  int default 10000 not null,
    tasks    int               null,
    username varchar(30)       null,
    password varchar(100)      null
);

create table loadertable
(
    ID        int auto_increment
        primary key,
    MaxWeight int        default 30    not null,
    Drunk     tinyint(1) default 0     not null,
    Fatigue   int        default 0     not null,
    Salary    int        default 10000 not null,
    username  varchar(30)              not null,
    password  varchar(100)             null
);



create table donetasks
(
    task_ID   int not null,
    loader_ID int not null,
    constraint l_ID
        foreign key (loader_ID) references loadertable (ID),
    constraint t_ID
        foreign key (task_ID) references tasks (ID)
);

create table tasks
(
    ID     int auto_increment
        primary key,
    name   varchar(20)    not null,
    weight int default 10 not null,
    done   int default 0  not null
);

