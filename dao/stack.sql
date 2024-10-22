use
    stack;
drop table if exists `profile`;
drop table if exists `names`;
drop table if exists `balances`;
drop table if exists `crawl`;
drop table if exists `cash_flow`;

CREATE TABLE `profile`
(
    `id`               bigint          NOT NULL AUTO_INCREMENT,
    `unique_id`        bigint unsigned not null default 0,
    `code`             varchar(20)     not null default '',
    `name`             varchar(200)    not null default '',
    `report_period`    varchar(30)     not null default '',
    `operate_in`       bigint          not null default 0,
    `operate_all_cost` bigint          not null default 0,
    `operate_cost`     bigint          not null default 0,
    `tax`              bigint          not null default 0,
    `sales_cost`       bigint          not null default 0,
    `manage_cost`      bigint          not null default 0,
    `rd_cost`          bigint          not null default 0,
    `financial_cost`   bigint          not null default 0,
    `invest`           bigint          not null default 0,
    `fair_in`          bigint          not null default 0,
    `net_profit`       bigint          not null default 0,
    `earn_per_share`   bigint          not null default 0,
    `created_at`       bigint          not null default 0,
    `updated_at`       bigint          not null default 0,
    PRIMARY KEY (`id`),
    unique index idx_unq (unique_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE TABLE `names`
(
    id           bigint       not null auto_increment primary key,
    `code`       varchar(10)  NOT NULL DEFAULT '',
    `name`       varchar(100) NOT NULL DEFAULT '',
    `price`      bigint       NOT NULL DEFAULT 0,
    `plate`      varchar(200) NOT NULL DEFAULT '',
    `industry`   varchar(200) NOT NULL DEFAULT '',
    `exchange`   varchar(200) NOT NULL DEFAULT '',
    `crawl_at`   bigint       NOT NULL DEFAULT 0,
    `created_at` bigint       NOT NULL DEFAULT 0,
    `updated_at` bigint       NOT NULL DEFAULT 0,
    UNIQUE index `idx_code` (code)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table balances
(
    id              bigint       not null auto_increment primary key,
    unique_id       bigint       not null default 0,
    code            varchar(10)  not null default '',
    name            varchar(200) not null default '',
    `report_period` varchar(30)  not null default '',
    money_funds     bigint       not null default 0,
    trans_finance   bigint       not null default 0,
    account_receive bigint       not null default 0,
    note_receive    bigint       not null default 0,
    account_pay     bigint       not null default 0,
    note_pay        bigint       not null default 0,
    assets          bigint       not null default 0,
    stock           bigint       not null default 0,
    construct       bigint       not null default 0,
    short_loan      bigint       not null default 0,
    long_loan       bigint       not null default 0,
    capital         bigint       not null default 0,
    `created_at`    bigint       NOT NULL DEFAULT 0,
    `updated_at`    bigint       NOT NULL DEFAULT 0,
    unique index `idx_code` (unique_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table crawl
(
    id           bigint      not null auto_increment primary key,
    unique_id    bigint      not null default 0,
    code         varchar(10) not null default '',
    `year`       varchar(50) not null default '',
    crawl_type   varchar(20) not null default '',
    `crawl_at`   bigint      NOT NULL DEFAULT 0,
    `created_at` bigint      NOT NULL DEFAULT 0,
    `updated_at` bigint      NOT NULL DEFAULT 0,
    unique index idx_code (unique_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

create table cash_flow
(
    id            bigint       not null auto_increment primary key,
    unique_id     bigint       not null default 0,
    name          varchar(300) not null default '',
    code          varchar(10)  not null default '',
    report_period varchar(30)  not null default '',
    sale_in       bigint       not null default 0,
    tax_in        bigint       not null default 0,
    sum_in        bigint       not null default 0,
    sale_out      bigint       not null default 0,
    emp_out       bigint       not null default 0,
    sum_out       bigint       not null default 0,
    netflow       bigint       not null default 0,
    created_at    bigint       NOT NULL DEFAULT 0,
    updated_at    bigint       NOT NULL DEFAULT 0,
    index idx_code (code),
    unique index unq_unique_id (unique_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;