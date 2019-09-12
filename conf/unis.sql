/*
 Navicat Premium Data Transfer

 Source Server         : 10.0.47.136
 Source Server Type    : MySQL
 Source Server Version : 50725
 Source Host           : 10.0.47.136:4000
 Source Schema         : unis

 Target Server Type    : MySQL
 Target Server Version : 50725
 File Encoding         : 65001

 Date: 12/09/2019 11:26:33
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
  INDEX `domain`(`project_id`, `host_domain`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 100 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
