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
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;


CREATE TABLE `names`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT,
    `code`       varchar(10)  NOT NULL DEFAULT '',
    `name`       varchar(100) NOT NULL DEFAULT '',
    `price`      bigint       NOT NULL DEFAULT 0,
    `plate`      varchar(200) NOT NULL DEFAULT '',
    `industry`   varchar(200) NOT NULL DEFAULT '',
    `exchange`   varchar(200) NOT NULL DEFAULT '',
    `crawl_at`   bigint       NOT NULL DEFAULT 0,
    `created_at` bigint       NOT NULL DEFAULT 0,
    `updated_at` bigint       NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_code` (`code`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1620
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci