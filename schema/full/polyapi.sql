/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_namespace`;
CREATE TABLE `api_namespace` (
    `id` 		    VARCHAR(64) 	NOT NULL COMMENT 'unique id',
    `owner` 	    VARCHAR(64) 	COMMENT 'owner id',
    `owner_name` 	VARCHAR(64) 	COMMENT 'owner name',
    `parent`        VARCHAR(320) 	NOT NULL DEFAULT '' COMMENT 'full namespace path, eg: /a/b/c',
    `namespace`     VARCHAR(64) 	NOT NULL  COMMENT 'global namespace, inmutable',
    `sub_count`     INT(10) DEFAULT 0 NOT NULL COMMENT 'count of sub namespace',
    `title` 	    VARCHAR(64) 	COMMENT 'alias of namespace, mutable',
    `desc` 		    TEXT,

    `access` 	    INT(11) 	    NOT NULL COMMENT 'privilege for public access, 1,2,4,8,16,32 CRUDGX',
    `active`        TINYINT         DEFAULT 1 COMMENT '1 ok, 0 disable',
    `valid`         TINYINT         DEFAULT 1 COMMENT '1 valid, 0 invalid',

    `create_at`     BIGINT(20) 		COMMENT 'create time',
    `update_at`     BIGINT(20) 		COMMENT 'update time',
    `delete_at`     BIGINT(20) 		COMMENT 'delete time',
    UNIQUE KEY `idx_global_name` (`parent`, `namespace`),
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_service`;
CREATE TABLE `api_service` (
    `id` 		    VARCHAR(64) 	COMMENT 'unique id',
    `owner` 	    VARCHAR(64) 	COMMENT 'owner id',
    `owner_name` 	VARCHAR(64) 	COMMENT 'owner name',
    `namespace`     VARCHAR(384) 	NOT NULL COMMENT 'full namespace path, eg: a/b/c',
    `name`	        VARCHAR(64) 	NOT NULL COMMENT 'service name, unique in namespace',
    `title` 	    VARCHAR(64) 	COMMENT 'alias of service, mutable',
    `desc` 		    TEXT,

    `access` 	    INT(11) 	 	NOT NULL COMMENT 'privilege for public access, 1,2,4,8,16,32 CRUDGX',
    `active`        TINYINT         DEFAULT 1 COMMENT '1 ok, 0 disable',
    `schema` 	    VARCHAR(16) 	NOT NULL COMMENT 'http/https',
    `host` 		    VARCHAR(128) 	NOT NULL COMMENT 'eg: api.xxx.com:8080',
    `auth_type`     VARCHAR(32)     NOT NULL COMMENT 'none/system/signature/cookie/oauth2...',
    `authorize`     TEXT            COMMENT 'JSON',

    `create_at`     BIGINT(20) 		COMMENT 'create time',
    `update_at`     BIGINT(20) 		COMMENT 'update time',
    `delete_at`     BIGINT(20) 		COMMENT 'delete time',
    UNIQUE KEY `idx_global_name` (`namespace`, `name`),
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

REPLACE  INTO `api_namespace`(`id`,`owner`,`owner_name`,`parent`,`namespace`,`title`,`desc`,`access`,`active`,`create_at`,`update_at`,`delete_at`) VALUES 
('1','system','系统','-','system','内部系统','系统自动注册的API',0,1,unix_timestamp(NOW())*1000,unix_timestamp(NOW())*1000,NULL),
('1-0','system','系统','/system','poly','内部聚合','内部生成的聚合API',0,1,unix_timestamp(NOW())*1000,unix_timestamp(NOW())*1000,NULL),
('1-1','system','系统','/system','faas','函数服务','通过faas注册的API',0,1,unix_timestamp(NOW())*1000,unix_timestamp(NOW())*1000,NULL),
('1-2','system','系统','/system','app','app根目录','app根目录',0,1,unix_timestamp(NOW())*1000,unix_timestamp(NOW())*1000,NULL),
('1-3','system','系统','/system','form','表单模型','通过form注册的API',0,1,unix_timestamp(NOW())*1000,unix_timestamp(NOW())*1000,NULL);

UPDATE `api_namespace` SET `sub_count`=0;
UPDATE `api_namespace` u, 
(SELECT `parent`, COUNT(1) cnt FROM `api_namespace` GROUP BY `parent`) t 
SET u.`sub_count`=t.`cnt` WHERE CONCAT(u.`parent`,'/',u.`namespace`)=t.`parent`;
UPDATE `api_namespace` u, 
(SELECT `parent`, COUNT(1) cnt FROM `api_namespace` GROUP BY `parent`) t 
SET u.`sub_count`=t.`cnt` WHERE CONCAT('/',u.`namespace`)=t.`parent`;

REPLACE  INTO `api_service`(`id`,`owner`,`owner_name`,`namespace`,`name`,`title`,`desc`,`access`,`active`,`schema`,`host`,`auth_type`,`authorize`,`create_at`,`update_at`,`delete_at`) VALUES 
('1','system','系统','/system/app','form','表单','表单接口',0,1,'http','form:8080','system',NULL,UNIX_TIMESTAMP(NOW())*1000,UNIX_TIMESTAMP(NOW())*1000,NULL),
('2','system','系统','/system/app','faas','函数服务','函数服务',0,1,'http','localhost:9999','none',NULL,UNIX_TIMESTAMP(NOW())*1000,UNIX_TIMESTAMP(NOW())*1000,NULL);

/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_raw`;
CREATE TABLE `api_raw` (
    `id`        VARCHAR(64)  COMMENT 'unique id',
    `owner`     VARCHAR(64)  COMMENT 'owner id',
    `owner_name`VARCHAR(64)  COMMENT 'owner name',
    `namespace` VARCHAR(384) NOT NULL COMMENT 'belong full namespace, eg: /a/b/c',
    `name`      VARCHAR(64) NOT NULL COMMENT 'unique name',
    `service`   VARCHAR(512) NOT NULL COMMENT 'belong service full path, eg: /a/b/c/servicesX',
    `title` 	VARCHAR(64)  COMMENT 'alias of name, mutable',
    `desc`      TEXT,
    `version`   VARCHAR(32),
    `path`      VARCHAR(512) NOT NULL COMMENT 'relative path, eg: /api/foo/bar',
    `url`       VARCHAR(512) NOT NULL COMMENT 'full path, eg: https://api.xxx.com/api/foo/bar',
    `action`    VARCHAR(64) DEFAULT '' COMMENT 'action on path',
    `method`    VARCHAR(16) NOT NULL COMMENT 'method',
    `content`   TEXT,
    `doc`   	TEXT COMMENT 'api doc',
    
    `access` 	INT(11) 	 NOT NULL COMMENT 'privilege for public access, 1,2,4,8,16,32 CRUDGX',
    `active`    TINYINT      DEFAULT 1 COMMENT '1 ok, 0 disable',
    `valid`         TINYINT         DEFAULT 1 COMMENT '1 valid, 0 invalid',
    
    `schema` 	VARCHAR(16) 	NOT NULL COMMENT 'from service, http/https',
    `host` 		VARCHAR(128) 	NOT NULL COMMENT 'eg: api.xxx.com:8080',
    `auth_type` VARCHAR(32)     NOT NULL COMMENT 'none/system/signature/cookie/oauth2...',
    
    `create_at` BIGINT(20) 	COMMENT 'create time',
    `update_at` BIGINT(20) 	COMMENT 'update time',
    `delete_at` BIGINT(20) 	COMMENT 'delete time',
    UNIQUE KEY `idx_global_name` (`namespace`,`name`),
    KEY `idx_service` (`service`),
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_poly`;
CREATE TABLE `api_poly` (
    `id`        VARCHAR(64)  COMMENT 'unique id',
    `owner`     VARCHAR(64)  COMMENT 'owner id',
    `owner_name`VARCHAR(64)  COMMENT 'owner name',
    `namespace` VARCHAR(384) NOT NULL COMMENT 'belong full namespace, eg: /a/b/c',
    `name`      VARCHAR(64) NOT NULL COMMENT 'name',
    `title` 	VARCHAR(64)  COMMENT 'alias of name, mutable',
    `desc`      TEXT,
    
    `access` 	INT(11) 	NOT NULL COMMENT 'privilege for public access, 1,2,4,8,16,32 CRUDGX',
    `active`    TINYINT     DEFAULT 1 COMMENT '1 ok, 0 disable',
    `valid`         TINYINT         DEFAULT 1 COMMENT '1 valid, 0 invalid',
    
    `method`    VARCHAR(16) NOT NULL COMMENT 'method',
    `arrange`   TEXT,
    `doc`       TEXT 		COMMENT 'api doc',
    `script`    TEXT,
    
    `create_at` BIGINT(20) 	COMMENT 'create time',
    `update_at` BIGINT(20) 	COMMENT 'update time',
    `build_at`  BIGINT(20) 	COMMENT 'build time',
    `delete_at` BIGINT(20) 	COMMENT 'delete time',
    UNIQUE KEY `idx_global_name` (`namespace`,`name`),
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8;

/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_raw_poly`;
Create TABLE `api_raw_poly` (
    `id`        VARCHAR(64)     NOT NULL COMMENT 'unique id',
    `raw_api`   VARCHAR(512)    NOT NULL COMMENT 'raw api full-path, eg: /a/b/c/rawApiName',
    `poly_api`  VARCHAR(512)    NOT NULL COMMENT 'poly api full-path, eg: /a/b/c/polyApiName',
    `delete_at` BIGINT(20)        COMMENT 'delete time',
    INDEX `idx_rawapi` (`raw_api`),
    INDEX `idx_polyapi` (`poly_api`),
    PRIMARY KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*--------------------------------------------------------------------------------------------------------------------------*/
DROP TABLE IF EXISTS `api_schema`;
CREATE TABLE `api_schema` (
    `ID`        VARCHAR(64)     COMMENT 'unique id',
    `namespace` VARCHAR(384)    NOT NULL COMMENT 'belong full namespace, eg: /a/b/c',
    `name`      VARCHAR(64)     NOT NULL COMMENT 'unique name',
    `title`     VARCHAR(64)     COMMENT 'alias of name',
    `desc`      TEXT,
    `schema`    TEXT            NOT NULL COMMENT 'api schema',
    `create_at` BIGINT(20)        COMMENT 'create time',
    `update_at` BIGINT(20)        COMMENT 'update time',
    `delete_at` BIGINT(20)        COMMENT 'delete time',
    UNIQUE KEY `idx_global_name` (`namespace`, `name`),
    PRIMARY KEY (`id`)
)ENGINE=INNODB DEFAULT CHARSET=utf8;
