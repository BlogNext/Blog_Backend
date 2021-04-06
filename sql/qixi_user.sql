/*
 Navicat Premium Data Transfer

 Source Server         : 晓琛-腾讯
 Source Server Type    : MySQL
 Source Server Version : 80023
 Source Host           : 154.8.142.48:3306
 Source Schema         : blog_xiaochen

 Target Server Type    : MySQL
 Target Server Version : 80023
 File Encoding         : 65001

 Date: 06/04/2021 11:56:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for attachment
-- ----------------------------
DROP TABLE IF EXISTS `attachment`;
CREATE TABLE `attachment`  (
  `id` int unsigned NOT NULL,
  `file_type` int unsigned NOT NULL COMMENT '文件类型 1: 图片, 2视频',
  `module` tinyint(0) NOT NULL COMMENT '功能模块 1：博客',
  `path` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '相对路径',
  `created_at` int(0) NOT NULL,
  `updated_at` int(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '附件表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for blog
-- ----------------------------
DROP TABLE IF EXISTS `blog`;
CREATE TABLE `blog`  (
  `id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL COMMENT '用户id，user表id',
  `doc_id` char(23) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '文档在es的唯一标识',
  `cover_plan_id` int(0) NOT NULL COMMENT '封面图',
  `blog_type_id` int unsigned NOT NULL COMMENT '博客类型id',
  `yuque_id` int(0) NOT NULL DEFAULT 0 COMMENT '语雀文档id',
  `yuque_slug` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '语雀文档路劲',
  `yuque_format` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '语雀文档格式',
  `yuque_lake` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语雀 lake格式文档',
  `yuque_public` tinyint(0) NOT NULL DEFAULT 0 COMMENT '语雀公开级别 0-私密,1公开',
  `title` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '博客标题',
  `abstract` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '博客摘要',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '博客内容，markdump',
  `browse_total` int unsigned NOT NULL COMMENT '浏览量',
  `created_at` int unsigned NOT NULL,
  `updated_at` int unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '博客' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for blog_type
-- ----------------------------
DROP TABLE IF EXISTS `blog_type`;
CREATE TABLE `blog_type`  (
  `id` int unsigned NOT NULL,
  `yuque_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '语雀知识库名',
  `yuque_id` int(0) NOT NULL COMMENT '语雀知识库id',
  `yuque_type` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '语雀知识库类型',
  `created_at` int unsigned NOT NULL,
  `updated_at` int unsigned NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_yuque_id`(`yuque_id`) USING BTREE COMMENT '语雀知识库唯一id'
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '博客类型(对饮语雀的知识库)' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for oauth_client
-- ----------------------------
DROP TABLE IF EXISTS `oauth_client`;
CREATE TABLE `oauth_client`  (
  `id` int unsigned NOT NULL,
  `client_id` char(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '格式: blog_1616579785;blog_时间戳',
  `client_secret` char(45) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '密码: blog_sha1函数机密,总长度45',
  `client_name` varchar(15) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用的名字',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uniq_clientId`(`client_id`) USING BTREE COMMENT 'clientId可以作为唯一标识记录'
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'oauth授权client模型' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` int unsigned NOT NULL,
  `nickname` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '昵称',
  `password` char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '密码  sha1加密',
  `created_at` int(0) NOT NULL COMMENT '创建时间',
  `updated_at` int(0) NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '用户表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for user_yuque
-- ----------------------------
DROP TABLE IF EXISTS `user_yuque`;
CREATE TABLE `user_yuque`  (
  `id` int unsigned NOT NULL COMMENT '语雀用户的id',
  `user_id` int unsigned NOT NULL COMMENT '用户id',
  `login` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语雀用户的login',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语雀用户的昵称',
  `avatar_url` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语雀用户头像',
  `description` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '语雀用户个性签名',
  `created_at` int(0) NOT NULL,
  `updated_at` int(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '语雀用户表，与user一对第一绑定' ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
