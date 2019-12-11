/*
 Navicat Premium Data Transfer

 Source Server         : 10.0.47.136
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : 10.0.42.26:4000
 Source Schema         : caddy

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 11/12/2019 17:42:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for custom_domain
-- ----------------------------
DROP TABLE IF EXISTS `custom_domain`;
CREATE TABLE `custom_domain`  (
  `project_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `host_domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `bucket_domain` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `tls_domain` blob NULL,
  `tls_domain_key` blob NULL,
  PRIMARY KEY (`project_id`, `host_domain`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for pipa
-- ----------------------------
DROP TABLE IF EXISTS `pipa`;
CREATE TABLE `pipa`  (
  `bucket_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `style_name` varchar(128) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `style_code` text CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`bucket_name`, `style_name`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_bin ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
