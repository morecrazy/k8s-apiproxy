CREATE DATABASE IF NOT EXISTS `kube_apiproxy`;

USE `kube_apiproxy`;

CREATE TABLE `slbs` (
	  `groupname` VARCHAR(128) NOT NULL,
	  `loadBalancerId`  VARCHAR(128) NOT NULL,
	  `loadBalancerName` VARCHAR(128),
	  `Ip` VARCHAR(128) NOT NULL,

	  `updated` datetime NOT NULL ,
	  `created` datetime NOT NULL,
	  PRIMARY KEY (`groupname`),
	  INDEX `idx_updated` (`updated`),
	  INDEX `idx_created` (`created`)
) ENGINE=Innodb DEFAULT CHARSET=utf8;
