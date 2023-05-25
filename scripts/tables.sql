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

INSERT INTO whitelist(`address`, `max_amount`, `proof`) VALUES('123', '234', '345');
