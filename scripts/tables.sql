DROP TABLE IF exists `process`;
CREATE TABLE IF NOT EXISTS `process` (
    `id`            BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `key`           VARCHAR(64)         NOT NULL                                                        COMMENT '任务的key',
    `value`         VARCHAR(256)        NOT NULL                                                        COMMENT '任务的进度',
    `create_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    UNIQUE KEY `uniq_address`(`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '任务进度';

DROP TABLE IF exists `events`;
CREATE TABLE IF NOT EXISTS `events` (
    `id`            BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `chain`         VARCHAR(64)         NOT NULL                                                        COMMENT '所在链',
    `block_hash`    CHAR(66)            NOT NULL                                                        COMMENT '所在区块HASH',
    `tx_index`      INT UNSIGNED        NOT NULL                                                        COMMENT '交易相对区块的位置',
    `log_index`     INT UNSIGNED        NOT NULL                                                        COMMENT '事件相对区块的位置',
    `contract`      CHAR(42)            NOT NULL                                                        COMMENT '发出事件的合约地址',
    `topics`        VARCHAR(330)        NOT NULL                                                        COMMENT '事件的topic,最多5个topic,多个topic以逗号连接',
    `data`          TEXT                                                                                COMMENT '事件的数据',
    `block_number`  INT UNSIGNED        NOT NULL                                                        COMMENT '所在区块块高',
    `tx_hash`       CHAR(66)            NOT NULL                                                        COMMENT '所在交易HASH',
    `removed`       SMALLINT            NOT NULL                                                        COMMENT '是否回滚 0表示未回滚 1表示已回滚',
    `create_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    UNIQUE KEY `uniq_key`(`chain`, `block_hash`, `log_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '链上事件';

DROP TABLE IF exists `user_chain_info`;
CREATE TABLE IF NOT EXISTS `user_chain_info` (
    `id`            BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `chain`         VARCHAR(64)         NOT NULL                                                        COMMENT '所在链',
    `address`       CHAR(42)            NOT NULL                                                        COMMENT '用户的钱包地址',
    `user_id`       BIGINT UNSIGNED                                                                     COMMENT '根据链上交易顺序生成的uid',
    `referral`      CHAR(42)                                                                            COMMENT '用户的邀请者',
    `referral_time` DATETIME                                                                            COMMENT '被邀请的时间',
    `invite_num`    INT UNSIGNED                                                                        COMMENT '用户的邀请人数量',
    `invite_reward` BIGINT UNSIGNED                                                                     COMMENT '用户的邀请奖励',
    `reward_num`    INT UNSIGNED                                                                        COMMENT '接收到的奖励次数',
    `create_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    UNIQUE KEY `uniq_address`(`chain`, `address`),
    UNIQUE KEY `uniq_user_id`(`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户的链上信息,由链上event决定';

DROP TABLE IF exists `referral_reward_info`;
CREATE TABLE IF NOT EXISTS `referral_reward_info` (
    `id`                BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `chain`             VARCHAR(64)         NOT NULL                                                        COMMENT '所在链',
    `address`           CHAR(66)            NOT NULL                                                        COMMENT '用户的钱包地址',
    `invite_address`    CHAR(66)            NOT NULL                                                        COMMENT '根据链上交易顺序生成的uid',
    `reward_amount`     BIGINT UNSIGNED     NOT NULL                                                        COMMENT '用户的邀请者',
    `reward_time`       DATETIME            NOT NULL                                                        COMMENT '被邀请的时间',
    `tx_hash`           CHAR(66)            NOT NULL                                                        COMMENT '被邀请的时间',
    `create_at`         DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`         DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    INDEX `idx_address`(`chain`, `address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户的邀请奖励';

DROP TABLE IF exists `price_info`;
CREATE TABLE IF NOT EXISTS `price_info` (
    `id`            BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `chain`         VARCHAR(64)         NOT NULL                                                        COMMENT '所在链',
    `price_type`    CHAR(64)            NOT NULL                                                        COMMENT '价格类型 1m 5m 15m 30m 1h 4h 1d 3d 7d 14d 1month 3month 6month',
    `price_key`     BIGINT UNSIGNED                                                                     COMMENT '价格记录的key',
    `price_open`    BIGINT UNSIGNED                                                                     COMMENT '',
    `price_high`    BIGINT UNSIGNED                                                                     COMMENT '',
    `price_low`     BIGINT UNSIGNED                                                                     COMMENT '',
    `price_close`   BIGINT UNSIGNED                                                                     COMMENT '',
    `create_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    UNIQUE KEY `uniq_key`(`chain`, `price_type`, `price_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '价格历史记录';

DROP TABLE IF exists `whitelist`;
CREATE TABLE IF NOT EXISTS `whitelist` (
    `id`            BIGINT UNSIGNED     NOT NULL PRIMARY KEY AUTO_INCREMENT                             COMMENT '自增id',
    `address`       VARCHAR(64)         NOT NULL                                                        COMMENT '用户钱包地址',
    `max_amount`    VARCHAR(256)        NOT NULL                                                        COMMENT '最大数量',
    `proof`         VARCHAR(2048)       NOT NULL                                                        COMMENT '证据',
    `create_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP                              COMMENT '创建时间',
    `update_at`     DATETIME            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  COMMENT '修改时间',
    UNIQUE KEY `uniq_address`(`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '白名单用户';
