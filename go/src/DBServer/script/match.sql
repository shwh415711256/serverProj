CREATE  TABLE  match_config(
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `match_type` int NOT NULL COMMENT '比赛类型',
  `ticket` int NOT NULL DEFAULT '0' COMMENT  '比赛门票',
  `losenum` int NOT NULL DEFAULT '0' COMMENT '失败扣除金币',
  `reward` varchar(2048) NOT NULL COMMENT '奖励信息,rankstart|rankend|rewardtype1-num,rewardtype2-num;rankstart|rankend|rewardtype1-num',
   PRIMARY  KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET="utf8mb4";
INSERT into match_config(`match_type`,`ticket`, `losenum`, `reward`)values(3, 2, 3,'1|1|1-3;2|2|1-0');
INSERT into match_config(`match_type`,`ticket`, `reward`)values(1, 40, '1|1|2-100;2|2|1-20;3|3|1-10');