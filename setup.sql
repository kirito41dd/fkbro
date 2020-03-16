DROP TABLE IF EXISTS `alert`;
CREATE TABLE `alert` (
    `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    `time_stamp` bigint UNSIGNED NOT NULL COMMENT '交易时间戳',
    `symbol` char(16) NOT NULL COMMENT '币种',
    `hash` char(64) NOT NULL COMMENT '交易hash',
    `amount` bigint DEFAULT 0,
    `amount_usd` bigint DEFAULT 0,
    `from_addr` char(64) DEFAULT NULL,
    `from_owner` char(32) DEFAULT NULL,
    `to_addr` char(64) DEFAULT NULL,
    `to_owner` char(32) DEFAULT NULL,

    PRIMARY KEY (`id`),
    UNIQUE KEY `UK_hash` (`hash`),
    INDEX IDX_time_stamp (time_stamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;