/*
 Navicat Premium Data Transfer

 Source Server         : 10.5.0.17
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : 10.5.0.17:4000
 Source Schema         : caddy

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 22/09/2019 21:47:17
*/

SET NAMES utf8;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for custom_domain
-- ----------------------------
DROP TABLE IF EXISTS `custom_domain`;
CREATE TABLE `custom_domain`  (
  `project_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `host_domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `bucket_domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`project_id`, `host_domain`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_bin;

SET FOREIGN_KEY_CHECKS = 1;
