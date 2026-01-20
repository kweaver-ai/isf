/*
MySQL: Database - authentication
*********************************************************************
*/

use authentication;

CREATE TABLE IF NOT EXISTS `t_session` (
    `f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT,
    `f_login_session_id` varchar(255) NOT NULL COMMENT 'session id',
    `f_subject` varchar(255) NOT NULL COMMENT '用户ID',
    `f_client_id` varchar(255) NOT NULL COMMENT '客户端ID',
    `f_exp` bigint(20) NOT NULL COMMENT 'Context到期时间戳',
    `f_session_access_token` text NOT NULL COMMENT 'Context信息',
    PRIMARY KEY (`f_primary_id`),
    UNIQUE KEY (`f_login_session_id`)
) ENGINE=InnoDB COMMENT='Context信息表';

CREATE TABLE IF NOT EXISTS `t_client_public` (
`id` varchar(255) NOT NULL COMMENT '客户端ID',
`client_name` text NOT NULL COMMENT '客户端名称',
`client_secret` text NOT NULL COMMENT '客户端密钥',
`redirect_uris` text NOT NULL COMMENT '客户端回调地址',
`grant_types` text NOT NULL COMMENT '客户端授权模式',
`response_types` text NOT NULL COMMENT '客户端接收响应类型',
`scope` text NOT NULL COMMENT '客户端申请权限范围',
`pk` int(10) unsigned NOT NULL AUTO_INCREMENT,
`post_logout_redirect_uris` text NOT NULL COMMENT '客户端登出成功地址',
`metadata` text NOT NULL COMMENT '客户端元数据',
`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '客户端创建时间',
PRIMARY KEY (`pk`),
UNIQUE KEY `hydra_client_idx_id_uq` (`id`)
) ENGINE=InnoDB COMMENT='公开注册客户端信息表';

CREATE TABLE IF NOT EXISTS `t_conf` (
  `f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `f_key` char(32) NOT NULL COMMENT '键',
  `f_value` varchar(1024) NOT NULL COMMENT '值',
  PRIMARY KEY (`f_primary_id`),
  UNIQUE KEY `uk_conf` (`f_key`)
) ENGINE=InnoDB COMMENT='认证配置表';

INSERT INTO t_conf(f_key,f_value) SELECT 'remember_for','2592000' FROM DUAL WHERE NOT EXISTS(SELECT f_value FROM t_conf WHERE f_key = 'remember_for');
INSERT INTO t_conf(f_key,f_value) SELECT 'remember_visible','true' FROM DUAL WHERE NOT EXISTS(SELECT f_value FROM t_conf WHERE f_key = 'remember_visible');

CREATE TABLE IF NOT EXISTS `t_access_token_perm` (
  `f_primary_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `f_app_id` char(36) NOT NULL COMMENT '应用账户id',
  `f_create_time` bigint(20) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`f_primary_id`),
  KEY `idx_f_app_id` (`f_app_id`)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS `t_outbox` (
    `f_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `f_business_type` tinyint(4) NOT NULL COMMENT '业务类型',
    `f_message` longtext NOT NULL COMMENT '消息内容，json格式字符串',
    `f_create_time` bigint(20) NOT NULL COMMENT '消息创建时间',
    PRIMARY KEY (`f_id`),
    KEY `idx_business_type_and_create_time` (`f_business_type`, `f_create_time`)
) ENGINE=InnoDB COMMENT='outbox信息表';

CREATE TABLE IF NOT EXISTS `t_outbox_lock` (
    `f_business_type` tinyint(4) NOT NULL COMMENT '业务类型',
    PRIMARY KEY (`f_business_type`)
) ENGINE=InnoDB COMMENT='outbox分布式锁表';

CREATE TABLE IF NOT EXISTS t_anonymous_sms_vcode (
    f_id char(26) NOT NULL COMMENT '验证码唯一标识',
    f_phone_number varchar(150) NOT NULL COMMENT '加密手机号',
    f_anonymity_id char(40) NOT NULL COMMENT '匿名账户id',
    f_content char(8) NOT NULL COMMENT '验证码内容',
    f_create_time timestamp NOT NULL COMMENT '创建时间',
    PRIMARY KEY (f_id),
    KEY idx_phone_number_anonymity_id (f_phone_number, f_anonymity_id),
    KEY idx_create_time (f_create_time)
) ENGINE=InnoDB COMMENT='匿名认证短信验证码存储表';

CREATE TABLE IF NOT EXISTS `t_distributed_lock` (
    `f_business_type` tinyint(4) NOT NULL COMMENT '业务类型',
    PRIMARY KEY (`f_business_type`)
) ENGINE=InnoDB COMMENT='分布式锁表';

CREATE TABLE IF NOT EXISTS `t_ticket` (
    `f_id` char(26) NOT NULL COMMENT '凭据唯一标识',
    `f_user_id` char(40) NOT NULL COMMENT '用户唯一标识',
    `f_client_id` varchar(255) NOT NULL COMMENT 'OAuth2客户端唯一标识',
    `f_create_time` bigint(10) NOT NULL COMMENT '凭据创建时间',
    PRIMARY KEY (`f_id`),
    KEY `idx_create_time` (`f_create_time`)
) ENGINE=InnoDB COMMENT='单点登录凭据表';

INSERT INTO t_outbox_lock(f_business_type) SELECT 1 FROM DUAL WHERE NOT EXISTS(SELECT f_business_type FROM t_outbox_lock WHERE f_business_type = 1);
INSERT INTO t_outbox_lock(f_business_type) SELECT 2 FROM DUAL WHERE NOT EXISTS(SELECT f_business_type FROM t_outbox_lock WHERE f_business_type = 2);
INSERT INTO t_outbox_lock(f_business_type) SELECT 3 FROM DUAL WHERE NOT EXISTS(SELECT f_business_type FROM t_outbox_lock WHERE f_business_type = 3);
INSERT INTO t_outbox_lock(f_business_type) SELECT 4 FROM DUAL WHERE NOT EXISTS(SELECT f_business_type FROM t_outbox_lock WHERE f_business_type = 4);

INSERT INTO t_distributed_lock(f_business_type) SELECT 1 FROM DUAL WHERE NOT EXISTS(SELECT f_business_type FROM t_distributed_lock WHERE f_business_type = 1);

INSERT INTO t_conf(f_key, f_value) SELECT 'anonymous_sms_expiration', '2' FROM DUAL WHERE NOT EXISTS(SELECT f_key FROM t_conf WHERE f_key = 'anonymous_sms_expiration');

CREATE TABLE IF NOT EXISTS t_outbox_unordered (
  id bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  f_message text NOT NULL COMMENT '信息',
  f_status tinyint(11) NOT NULL DEFAULT 0 COMMENT '状态(0 未开始,1 处理中)',
  f_created_at bigint(20) NOT NULL COMMENT '创建时间',
  f_updated_at bigint(20) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (id),
  KEY idx_f_status(f_status),
  KEY idx_f_created_at(f_created_at),
  KEY idx_f_updated_at(f_updated_at)
)ENGINE = InnoDB COMMENT='无序outbox信息表';
