-- 1.1 查询条件包含or，可能导致索引失效
CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` int(11) NOT NULL,
    `age` int(11) NOT NULL,
    `name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_userId` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- 走索引
EXPLAIN SELECT * FROM `user` where userId = 1;
-- 索引失效
EXPLAIN SELECT * FROM `user` where userId = 1 or age = 10;

-- 走索引
EXPLAIN SELECT * FROM `user` where userId = 1 or userId = 10;

-------------------------------------------------------------
-- 1.2  如果字段类型是字符串，where是一定用引号括起来，否则索引会失效。
DROP Table `user`;

CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` varchar(32) NOT NULL,
    `name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_userId` (`userId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
-- 不走索引
EXPLAIN SELECT * FROM `user` WHERE userId = 123;
-- 走索引
EXPLAIN SELECT * FROM `user` WHERE userId = "123";

-------------------------------------------------------------
-- 1.3 like通配符可能导致索引失效
-- 1.  like查询以%开头。索引失效
EXPLAIN SELECT * FROM `user` WHERE userId LIKE '%123';
-- 2.  把%放后面，发现索引还是正常走的
EXPLAIN SELECT * FROM `user` WHERE userId LIKE '123';
-- 3.  把%加回来，改为只查索引的字段(索引覆盖)，发现还是走的索引
EXPLAIN SELECT userId,id FROM `user` WHERE userId LIKE '%123%';

-------------------------------------------------------------
-- 1.4 联合索引，查询是的条件列不是联合索引中的第一个列，索引失效
ALTER TABLE `user` DROP INDEX `idx_userId`;

ALTER TABLE `user` MODIFY COLUMN `userId` INT(11);

ALTER TABLE `user` ADD COLUMN age INT(11) NOT NULL;

CREATE INDEX `idx_userId_age` ON `user` (`userId`,`age`) USING BTREE;

-- 在联合索引中，查询条件满足最左匹配原则时，索引是正常生效的
EXPLAIN SELECT * FROM `user` WHERE userId=10 and age=10;

-------------------------------------------------------------
-- 1.5 在索引列上使用mysql的内置函数，索引失效
ALTER TABLE `user` DROP INDEX  `idx_userId_age`;

ALTER TABLE `user` MODIFY `userId` varchar(32) NOT NULL;

ALTER TABLE `user` ADD COLUMN `loginTime` datetime NOT NULL;

CREATE INDEX `idx_userId` ON `user` (`userId`) USING BTREE;

CREATE INDEX `idx_login_time` ON `user` (`loginTime`) USING BTREE;

SHOW CREATE TABLE `user`;

-- 虽然loginTime加了索引，但是因为使用了mysql的内置函数Date_add()，索引直接失效
EXPLAIN SELECT * from `user` WHERE DATE_ADD(loginTime,INTERVAL 1 DAY) <= 7;

-------------------------------------------------------------
-- 1.6 对索引列运算(如：+ — * /)，索引失效。
SHOW CREATE TABLE `user`;
ALTER Table `user` DROP KEY `idx_login_time`;

ALTER Table `user` DROP COLUMN `loginTime`;
ALTER Table `user` DROP COLUMN `name`;
ALTER Table `user` DROP KEY `idx_userId`;
CREATE INDEX `idx_age` ON `user` (`age`) USING BTREE;
-- 虽然加了索引，但是因为它进行运算，索引失效了
EXPLAIN SELECT * FROM `user` WHERE age-1 =10;

-------------------------------------------------------------
-- 1.7 索引字段上使用(!= 或 <>, not in)时，可能会导致索引失效。
-- 虽然age加了索引，但是使用了!= 或者<>,not in这些时，索引失效了
EXPLAIN SELECT * FROM `user1` WHERE age <> 10;


-------------------------------------------------------------
-- 1.8 索引字段上使用is null, is not null, 可能导致索引失效。
SHOW CREATE TABLE `user`;
ALTER TABLE `user` add COLUMN `card` varchar(255) DEFAULT NULL;
ALTER TABLE `user` add COLUMN `name` varchar(255) DEFAULT NULL;
CREATE INDEX `idx_card` ON `user` (`card`) USING BTREE;
CREATE INDEX `idx_name` ON `user` (`name`) USING BTREE;

ALTER TABLE `user` CHANGE age age int(32) DEFAULT NULL;

ALTER TABLE `user` CHANGE card card varchar(255) DEFAULT NULL;
-- 单个name字段加上索引，并查询name为非空的语句，其实会走索引的
EXPLAIN SELECT * FROM `user` WHERE name IS NOT NULL;
-- 单个name字段加上索引，并查询card为非空的语句，其实会走索引的
EXPLAIN SELECT * FROM `user` WHERE card IS NOT NULL;
-- 但是它俩用or连接起来，索引就失效了
EXPLAIN SELECT * FROM `user` WHERE card IS NOT NULL or name IS NOT NULL;

-------------------------------------------------------------
-- 1.9 左连接查询或者右连接查询查询关联的字段编码格式不一样，可能导致索引失效。
DROP Table user;

CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL,
  `age` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

CREATE TABLE `user_job` (
  `id` int(11) NOT NULL,
  `userId` int(11) NOT NULL,
  `job` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 执行左外连接查询，user_job表还是走的全表扫描
EXPLAIN SELECT u.name,j.job FROM `user` u LEFT JOIN user_job j on u.name = j.name;