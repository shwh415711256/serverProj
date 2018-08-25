CREATE  TABLE  version_config(
  `id` bigint(20) unsigned NOT NULL auto_increment,
  `key` varchar(128) NOT NULL DEFAULT '' COMMENT 'key:gameid',
  `ignore_flag` tinyint(1) NOT NULL DEFAULT '0',
   PRIMARY  KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET="utf8mb4";

INSERT into version_config(`key`, `ignore_flag`)VALUES ('1', '1');
INSERT into version_config(`key`, `ignore_flag`)VALUES ('2', '1');
INSERT into version_config(`key`, `ignore_flag`)VALUES ('3', '0');
INSERT into version_config(`key`, `ignore_flag`)VALUES ('4', '0');
INSERT into version_config(`key`, `ignore_flag`)VALUES ('5', '0');
INSERT into version_config(`key`, `ignore_flag`)VALUES ('6', '0');