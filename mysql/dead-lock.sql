use fly;
-- 1.  数据库隔离级别
SELECT @@transaction_isolation;
-- 2.  自动提交关闭
set autocommit=0;
SELECT @@autocommit;

CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `balance` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- 插入数据：
INSERT INTO `account` VALUES (1,"Eason",100), (2,"Wei",100);

-- 开启两个事务
-- 事务A开始
begin;
UPDATE account SET balance = 1000 WHERE name = "Wei";

-- 事务B开始：
BEGIN;
UPDATE `account` SET balance = 1000 WHERE name = "Eason";

-- 事务A插入数据，陷入阻塞
INSERT INTO `account` VALUES(NULL,'Jay',100);

-- 查看锁情况：
SELECT * FROM performance_schema.data_locks;
SELECT * FROM performance_schema.data_lock_waits;

-- 事务B执行插入操作，插入成功，同时事务A的插入有阻塞变为死锁error
INSERT INTO `account` VALUES(null,"Yan",100);

SELECT * FROM `account`;


SELECT * FROM information_schema.INNODB_TRX;

-- 记录锁是最简单的行锁，仅仅锁住一行
SELECT * FROM account WHERE balance = 100 FOR UPDATE;

show engine innodb status;