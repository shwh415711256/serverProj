CREATE  TABLE  sign_reward_config(
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `type` tinyint unsigned NOT NULL COMMENT '签到类型',
  `reward` varchar(2048) NOT NULL DEFAULT '' COMMENT '签到奖励数据 1|1-10,2-100;2|1-20',
  PRIMARY  KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET="utf8mb4";
insert into sign_reward_config(`type`, `reward`)VALUES (1, '1|1-50;2|1-50;3|1-50;4|1-50;5|1-50;6|1-50;6|1-50');


CREATE  TABLE  sign_info_5(
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `sign_type` tinyint unsigned NOT NULL COMMENT '签到类型',
  `openid` varchar(128) NOT NULL COMMENT '微信用户openid',
  `sign_day` tinyint NOT NULL DEFAULT '0' COMMENT '领取天数',
  `today_sign` tinyint(1) NOT NULL DEFAULT '0' COMMENT '今天是否已签到',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY  KEY (`id`),
  KEY `o_index` (`openid`) USING BTREE
)ENGINE=InnoDB DEFAULT CHARSET="utf8mb4";