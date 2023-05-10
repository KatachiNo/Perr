create table "Products"
(
    category          integer not null,
    picture_address   varchar(100),
    id                serial
        primary key
        unique,
    quantity_of_goods integer not null,
    last_price        numeric,
    product_name      varchar(100),
    available_status  varchar(50)
);

comment on column "Products".category is 'number of category ';

comment on column "Products".picture_address is 'Address of picture in file system';

comment on column "Products".product_name is 'Name of product';

comment on column "Products".available_status is 'available/not available/in stock - vars';

alter table "Products"
    owner to postgres;

create table "CategoryTable"
(
    categoryid   serial
        unique,
    categoryname varchar(100) not null,
    id           serial
        constraint "CategoryTable_pk"
            primary key
);

alter table "CategoryTable"
    owner to postgres;

create table "ProductPriceStory"
(
    id      serial
        constraint "ProductPrice_pkey"
            primary key
        unique,
    "Price" numeric   not null,
    "Date"  timestamp not null
);

alter table "ProductPriceStory"
    owner to postgres;

create table "Users"
(
    id                   serial
        constraint "Logins_pk"
            primary key
        unique,
    login                varchar(100)  not null
        unique,
    "passwordHash"       varchar(1000) not null,
    "categoryOfUser"     varchar(100)  not null,
    "dateOfRegistration" timestamp     not null,
    salt                 varchar(1000) not null,
    algorithm            varchar(200)  not null
);

alter table "Users"
    owner to postgres;

create table "UserData"
(
    id           serial
        constraint id
            primary key
        unique,
    email        varchar(50),
    phone_number varchar(30),
    country      varchar(30),
    city         varchar(30),
    index        varchar(30),
    street       varchar(30),
    number_house varchar(10),
    note         varchar(100),
    first_name   varchar(50),
    middle_name  varchar(50),
    last_name    varchar(50)
);

alter table "UserData"
    owner to postgres;

create table "Orders"
(
    "orderId"            serial
        constraint "Orders_pk"
            primary key
        unique,
    user_id              integer   not null,
    data_of_order        timestamp not null,
    ordered_products_ids integer[] not null,
    final_price          numeric   not null,
    delivery_status      varchar(50)
);

comment on column "Orders"."orderId" is 'order id';

comment on column "Orders".user_id is 'who ordered';

comment on column "Orders".ordered_products_ids is 'ids of products which was ordered';

alter table "Orders"
    owner to postgres;
