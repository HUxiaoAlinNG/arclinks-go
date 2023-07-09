```sql
DROP TABLE IF EXISTS `onetime_password`;
CREATE TABLE `onetime_password` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `open_key_id` int(11) NOT NULL DEFAULT '0' COMMENT 'open_key表id',
  `verify` tinyint(2) NOT NULL DEFAULT '0' COMMENT '是否已验证 0否1是',
  `password` varchar(100) NOT NULL DEFAULT '' COMMENT '一次性密码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_openkeyid` (`open_key_id`) USING BTREE,
  KEY `idx_password` (`password`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COMMENT='一次性密码表';

-- ----------------------------
-- Table structure for open_key
-- ----------------------------
DROP TABLE IF EXISTS `open_key`;
CREATE TABLE `open_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `key` varchar(200) NOT NULL DEFAULT '' COMMENT 'openKey',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_key` (`key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8mb4 COMMENT='openkey表';

SET FOREIGN_KEY_CHECKS = 1;
```