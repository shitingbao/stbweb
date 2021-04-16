CREATE TABLE `task` (
  `task_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '任务名称',
  `sys` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '对应系统',
  `spec` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '执行定时cron字符串',
  `is_open` int NOT NULL DEFAULT '0' COMMENT '0关闭，1开启，2启动中,3关闭中',
  `version` float(10,2) NOT NULL DEFAULT '1.00' COMMENT '版本，沿用Linux',
  PRIMARY KEY (`task_name`,`sys`) USING BTREE
) 

CREATE TABLE `task_err` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_name` varchar(255) DEFAULT NULL,
  `sys` varchar(255) DEFAULT NULL,
  `runtime_error` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) 

CREATE TABLE `task_history` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `task_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '任务名称',
  `sys` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '对应系统',
  `spec` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '',
  `err` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '错误信息',
  `start_time` datetime NOT NULL DEFAULT '0000-01-01 00:00:00' COMMENT 'start时间',
  `end_time` datetime NOT NULL DEFAULT '0000-01-01 00:00:00' COMMENT '结束时间',
  `is_complete` int NOT NULL DEFAULT '0' COMMENT '0完成，1失败',
  `runtime_error` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '执行过程err',
  PRIMARY KEY (`id`) USING BTREE
)