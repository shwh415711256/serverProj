CREATE  TABLE  game_his_5(
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `openid` varchar(128) NOT NULL COMMENT '微信用户openid',
  `hisdata` varchar(4096) NOT NULL COMMENT '历史数据,json字串',
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   PRIMARY  KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET="utf8mb4";

insert game_his_5(`openid`, `hisdata`) values('12345678', '{"Data":[{"Time":2000,"Score":3},{"Time":4000,"Score":5},{"Time":7000,"Score":8}]}');
insert game_his_5(`openid`, `hisdata`) values('1234567', '{"Data":[{"Time":2121,"Score":1},{"Time":30000,"Score":2},{"Time":6000,"Score":4},{"Time":10000,"Score":4}]}');
insert game_his_5(`openid`, `hisdata`) values('123456', '{"Data":[{"Time":2000,"Score":1},{"Time":3000,"Score":2},{"Time":5000,"Score":5},{"Time":10000,"Score":8}]}');