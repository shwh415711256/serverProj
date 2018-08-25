drop table user_info_tbl_1;
drop table  user_info_tbl_2;
drop table  user_info_tbl_3;
drop table  user_info_tbl_4;
drop table user_info_tbl_5;
drop table user_info_tbl_6;
CREATE TABLE user_info_tbl_1(
  `uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `openid` varchar(128) NOT NULL COMMENT '微信用户openid',
  `gold` bigint(20) NOT NULL DEFAULT '50' COMMENT '金币',
  `money` bigint(20) NOT NULL DEFAULT '0' COMMENT '现金',
  `score` bigint(20) NOT NULL DEFAULT '0' COMMENT '分数',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`uid`),
  unique KEY `u_openid` (`openid`) USING BTREE
)ENGINE=InnoDB DEFAULT charset="utf8mb4";
create table user_info_tbl_2 like user_info_tbl_1;
create table user_info_tbl_3 like user_info_tbl_1;
create table user_info_tbl_4 like user_info_tbl_1;
create table user_info_tbl_5 like user_info_tbl_1;
create table user_info_tbl_6 like user_info_tbl_1;
insert into user_info_tbl_5(`openid`) values("12345678");
insert into user_info_tbl_5(`openid`) values("1234567");
insert into user_info_tbl_5(`openid`) values("123456");
insert into user_info_tbl_5(`openid`) values("testOpenId");

drop table wechat_user_tbl_1;
drop table wechat_user_tbl_2;
drop table wechat_user_tbl_3;
drop table wechat_user_tbl_4;
drop table wechat_user_tbl_5;
drop table wechat_user_tbl_6;
CREATE TABLE wechat_user_tbl_1(
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `openid` varchar(128) NOT NULL COMMENT '微信用户openid',
    `nick_name` varchar(128) NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar_url` varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
    `gender` char(1) NOT NULL DEFAULT '' COMMENT '性别',
    `city` varchar(40) NOT NULL DEFAULT '' COMMENT '城市',
    `province` varchar(40) NOT NULL DEFAULT '' COMMENT '省份',
    `country` varchar(40) NOT NULL DEFAULT '' COMMENT '国家',
    `language` varchar(40) NOT NULL DEFAULT '' COMMENT '语言',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`),
    unique KEY `u_openid` (`openid`) USING BTREE
)ENGINE=InnoDB DEFAULT charset="utf8mb4";
create table wechat_user_tbl_2 like wechat_user_tbl_1;
create table wechat_user_tbl_3 like wechat_user_tbl_1;
create table wechat_user_tbl_4 like wechat_user_tbl_1;
create table wechat_user_tbl_5 like wechat_user_tbl_1;
create table wechat_user_tbl_6 like wechat_user_tbl_1;
insert into wechat_user_tbl_5(`openid`, `nick_name`) values("12345678", "testssssss");
insert into wechat_user_tbl_5(`openid`, `nick_name`) values("1234567", "test2222222222");
insert into wechat_user_tbl_5(`openid`, `nick_name`) values("123456", "test3333333333");
insert into wechat_user_tbl_5(`openid`, `nick_name`) values("testOpenId2", "test55555555");

/*
CREATE TABLE user_login_his_1(
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `openid` varchar(128) NOT NULL COMMENT '微信用户openid',
  `num` varchar(20) NOT NULL DEFAULT '1' COMMENT '数量',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY(`id`),
  unique KEY `u_openid` (`openid`) USING BTREE
)ENGINE=InnoDB DEFAULT charset="utf8mb4";
create table user_login_his_2 like user_login_his_1;
create table user_login_his_3 like user_login_his_1;
create table user_login_his_4 like user_login_his_1;
create table user_login_his_5 like user_login_his_1;
create table user_login_his_6 like user_login_his_1;
*/

CREATE TABLE qd_login_info(
    `chanid` varchar(10) NOT NULL COMMENT '渠道',
    `gameid` varchar(10) NOT NULL COMMENT '游戏id',
    `num` varchar(64) NOT NULL DEFAULT '0' COMMENT '人数',
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    unique KEY `cg_index` (`chanid`, `gameid`) USING BTREE
)ENGINE=InnoDB DEFAULT charset="utf8mb4";