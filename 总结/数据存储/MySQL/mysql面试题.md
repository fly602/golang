# 1、 MySQL索引使用有哪些注意事项呢？
```
可以从三个维度回答这个问题：索引哪些情况会失效，索引不适合哪些场景，索引规则。
```
## 索引哪些情况会失效
### 1.1 查询条件包含or，可能导致索引失效
新建一个user表，它有一个普通索引userid，结构如下：
```
CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` int(11) NOT NULL,
    `age` int(11) NOT NULL,
    `name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_userId` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```
1.  执行一条查询sql：“EXPLAIN SELECT * FROM `user` where userId = 1;” 它是会走索引的。
2.  把or条件+没有索引的age加上，并不会走索引：“EXPLAIN SELECT * FROM `user` where userId = 1 or age = 10;”

####    分析&结论：
-   对于or+没有索引的age这种情况，假设它走了userid的索引，但是走到age查询条件时，它还得全表扫描，也就是需要三步过程：全表扫描+索引扫描+合并
-   如果它一开始就走全表扫描，直接一遍扫描就完事。
-   mysql是有优化器的，处于效率与成本，遇到or条件，索引可能失效，看起来也合情合理。
-   如果or条件的列都加了索引，索引可能会走的。例如：“EXPLAIN SELECT * FROM `user` where userId = 1 or userId = 10;”

### 1.2 如果字段类型是字符串，where是一定用引号括起来，否则索引会失效。
新建一个user表：
```
CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` varchar(32) NOT NULL,
    `name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_userId` (`userId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
1.  userid为字符串类型，是B+树的普通索引，如果查询条件穿了一个数字过去，它是不走索引的：“EXPLAIN SELECT * FROM `user` WHERE userId = 123;”
2.  如果给数字加上"",也就是传一个字符串，就走索引：“EXPLAIN SELECT * FROM `user` WHERE userId = "123";”

####    分析与结论：
为什么第一条语句未加引号就不走索引了呢？这是因为不加单引号时，是字符串跟数字的比较，他们类型不匹配，Mysql会做隐式的类型转换，把他们转换成浮点数在做比较。

### 1.3 like通配符可能导致索引失效
并不是用了like通配符，索引一定失效，而是like查询是以%开头，才会导致索引失效。
表结构：
```
CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `userId` varchar(32) NOT NULL,
    `name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_userId` (`userId`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
1.  like查询以%开头。索引失效：“EXPLAIN SELECT * FROM `user` WHERE userId LIKE '%123';”
2.  把%放后面，发现索引还是正常走的：“EXPLAIN SELECT * FROM `user` WHERE userId LIKE '123';”
3.  把%加回来，改为只查索引的字段(索引覆盖)，发现还是走的索引：“EXPLAIN SELECT userId,id FROM `user` WHERE userId LIKE '%123%';”

####    分析&结论：
like查询以%开头，会导致索引失效。可以有两种方式优化：
-   使用索引覆盖
-   把%放后面

### 1.4 联合索引，查询是的条件列不是联合索引中的第一个列，索引失效
表结构：（有一个联合索引idx_userid_age,userId在前，age在后）
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` int(11) NOT NULL,
  `age` int(11) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_userid_age` (`userId`,`age`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
在联合索引中，查询条件满足最左匹配原则时，索引是正常生效的：“EXPLAIN SELECT * FROM `user` WHERE userId=10 and age=10;”

####    分析&结论
-   当我们创建一个联合索引的时候。如（k1,k2,k3），相当于创建了（k1）、（k1,k2）和（k1,k2,k3）三个索引，这就是最左匹配原则。
-   联合索引不满足最左匹配原则，索引一般会失效，但是这个还是跟mysql优化器有关的。

### 1.5 在索引列上使用mysql的内置函数，索引失效。
表结构：
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` varchar(32) NOT NULL,
  `loginTime` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_userId` (`userId`) USING BTREE,
  KEY `idx_login_time` (`loginTime`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
虽然loginTime加了索引，但是因为使用了mysql的内置函数Date_add()，索引直接失效：“EXPLAIN SELECT * from `user` WHERE DATE_ADD(loginTime,INTERVAL 1 DAY) <= 7;”

### 1.6 对索引列运算(如：+ — * /)，索引失效。
表结构：
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` varchar(32) NOT NULL,
  `age` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_age` (`age`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
虽然加了索引，但是因为它进行运算，索引失效了：“EXPLAIN SELECT * FROM `user` WHERE age =10;”

### 1.7 索引字段上使用(!= 或 <>, not in)时，可能会导致索引失效。（此问题mysql高版本可能不存在）
表结构：
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` int(11) NOT NULL,
  `age` int(11) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_age` (`age`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```
虽然age加了索引，但是使用了!= 或者<>,not in这些时，索引失效了：“EXPLAIN SELECT * FROM `user1` WHERE name <> "10";”

### 1.8 索引字段上使用is null, is not null, 可能导致索引失效。
表结构：
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `card` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE,
  KEY `idx_card` (`card`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
```

单个name字段加上索引，并查询name为非空的语句，其实会走索引的：“EXPLAIN SELECT * FROM `user` WHERE name IS NOT NULL;”
单个name字段加上索引，并查询card为非空的语句，其实会走索引的：“EXPLAIN SELECT * FROM `user` WHERE card IS NOT NULL;”
但是它俩用or连接起来，索引就失效了：“EXPLAIN SELECT * FROM `user` WHERE card IS NOT NULL or name IS NOT NULL;”

### 1.9 左连接查询或者右连接查询查询关联的字段编码格式不一样，可能导致索引失效。
新建两个表，一个user，一个user_job:
```
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

```
user表的name字段编码是utf8mb4,而user_job表的name字段编码是utf8.
执行左外连接查询，user_job表还是走的全表扫描：“EXPLAIN SELECT u.name,j.job FROM `user` u LEFT JOIN user_job j on u.name = j.name;”

### 1.10  mysql估计使用全表扫描比使用索引快，则不适用索引。
- 当表的索引被查询，会使用最好的索引，除非优化器使用全表扫描更有效。优化器优化成全表扫描取决于使用最好的索引查询出来的数据是否超过全表30%的数据。
- 不要给“性别”等增加索引，如果某个数据列里包含了均是“0/1”或“Y/N”等值，即包含着许多重复的值，就算为它建立了索引，索引的效果不会太好，还可能导致全表扫描。

### 索引潜规则
- 覆盖索引
- 回表
- 索引数据结构（B+树）
- 最左前缀原则
- 索引下推

# Mysql死锁问题，如何解决死锁问题。
排查死锁的一般方法：
- 查看死锁日志show engine innodb status;
- 找出死锁sql
- 模拟死锁案发
- 分析死锁日志
- 分析死锁结果

##  环境准备
1.  数据库隔离级别：
mysql version < v8:
select @@tx_isolation;

mysql version > v8:
select @@transaction_isolation;

2.  自动提交关闭：
set autocommit=0;

3.  表结构：
```
//id是自增主键，name是非唯一索引，balance普通字段
CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `balance` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

```
创建两个事务进行模拟操作：
1.  事务A执行更新操作，更新成功
```
begin;
UPDATE account SET balance = 1000 WHERE name = "Wei";
```
2.  事务B执行更行操作，更行成功
```
BEGIN;
UPDATE `account` SET balance = 1000 WHERE name = "Eason";
```
3.  事务A执行插入操作，陷入阻塞
```
INSERT INTO `account` VALUES(NULL,'Jay',100);
```

4.  查询锁情况：
- mysql v8版本以后：
```
// 查询死锁表
SELECT * FROM performance_schema.data_locks;
// 查询死锁等待时间
SELECT * FROM performance_schema.data_lock_waits;
```

- mysql v8版本以前：
```
// 查询死锁表
SELECT * FROM INFORMATION_SCHEMA.INNODB_LOCKS;
// 查询死锁等待时间
SELECT * FROM information_schema.INNODB_LOCK_waits;
```
5.  事务B执行插入操作，插入成功，同时事务A的插入有阻塞变为死锁error。
```
INSERT INTO `account` VALUES(null,"Yan",100);
```

##  mysql 锁介绍：
在分析死锁日志前，先做一下锁介绍：
- 加锁机制：乐观锁、悲观锁
- 锁粒度：表锁、页锁、行锁
- 兼容性：共享锁、排他锁
- 锁模式：记录锁、gap锁、next-key锁、意向锁、插入意向锁

主要介绍一下兼容性锁和锁模式类型的锁：
### 共享锁和排他锁：
InnoDB 实现了标准的行级锁，包括两种：共享锁（简称s锁）、排他锁（x锁）。
- 共享锁（s锁）：允许持锁事务读取一行。
- 排他锁（X锁）：允许持锁事务更新和删除一行。

如果事务T1持有行R的S锁，那么另一个事务T2请求R的锁时，会做如下处理：
- T2请求S锁立即被允许，结果T1T2都持有R行的S锁
- T2请求X锁不能立即允许

如果T1 持有R的X锁，那么T2请求R的X、S锁都不能被立即允许，T2必须等待T1释放X锁才可以，因为X锁与任何锁都不兼容。

### 意向锁：
- 意向共享锁（IS锁）：事务想要获得一张表中某几行的共享锁，事务在请求S锁前，要先获取IS锁
- 意向排他锁（IX锁）：事务想要获取一张表中某几行的排他锁，事务在请求X锁前，要先获取IX锁

比如：事务1在表1加上了S锁后，事务2想要更改某行记录，需要添加IX锁，由于不兼容，所以需要等待S锁释放；如果事务1在表1加上了IS锁，事务2添加的IX锁与IS锁兼容，就可以操作，这就实现了更细粒度的加锁。

InnoDB存储引擎中锁的兼容性如下表：
+-------+------------------------------+
| 兼容性|  IS    |IX    |S      |X      |
+--------------------------------------+
| IS    | 兼容   |兼容  |兼容    |不兼容  |
+--------------------------------------+
| IX    | 兼容   |兼容  |不兼容   |不兼容 |
+--------------------------------------+
| S     | 兼容   |不兼容|兼容    |不兼容  |
+--------------------------------------+
| X     | 不兼容 |不兼容|不兼容   |不兼容  |
+-------+------------------------------+

### 记录锁（Record Locks）
- 记录锁是最简单的行锁，仅仅锁住一行。如：SELECT C1 from T where C1=10 for update；
- 记录锁永远是在索引上加的，实际一个表没有索引，InnoDB也会隐式的创建一个索引，并使用这个索引试试记录锁。
- 会阻塞其他事务对其插入、更新、删除。

### 间隙锁：
- 间隙锁是一种加载两个索引之间的锁，或者加载第一个索引之前，或者最后一个索引之后的间隙。
- 使用间隙锁锁住的是一个区间，而不仅仅是这个区间中的每一条数据。
- 间隙锁只阻止其他事务插入到间隙中，他们不阻止其他事务在同一个间隙上获得间隙锁所以gap x lock 和gap s lock有相同的作用。

### Next-key Locks
- Next-key锁是记录锁和间隙锁的组合，他指的是加在某条记录以及这条记录前面间隙上的锁。

### 插入意向锁（Insert Intention）
- 插入意向锁是在插入一行记录操作之前设置的一种间隙锁，这个锁释放了一种插入方式的信号，亦即多个事务在相同的索引间隙插入时如果不是插入间隙中相同的位置就不需要互相等待。
- 假设有索引值4、7，几个不同的事务准备插入5、6，每个锁都在获得插入行的独占锁之前用插入意向锁各自锁住4、7之间的间隙，但是不阻塞对方因为插入行不冲突。

### 锁模式兼容矩阵(横向是已持有锁，纵向是正在请求的锁)
+----------+----------+----------+----------+----------+
| 兼容性    | Gap      | Insert   | Record   | Next Key |
|          |          | Intention|          |          |
+----------+-------------------------------------------+
| Gap      | 兼容     | 兼容     | 兼容       | 兼容      |
+------------------------------------------------------+
| Insert   | 冲突     | 兼容     | 兼容       | 冲突      |
| Intention|          |          |          |          |
+------------------------------------------------------+
| Record   | 兼容     | 兼容     | 冲突       | 冲突      |
+------------------------------------------------------+
| Next Key | 兼容     | 兼容     | 冲突       | 冲突      |
+----------+----------+----------+----------+----------+


##  如何读懂死锁日志
### show engine innodb status

## 死锁死循环四要素
- 互斥条件：指进程对所分配到的资源进行排他性使用，即在一段时间内某资源只能由一个进程占用。如果此时还有其他进程请求资源，则请求者只能等待，直至战友资源的进程释放资源。
- 请求和保持条件：指进程已经保持只有一个资源，单有提出了新的资源请求，而该资源已经被其他进程战友，此时请求进程阻塞，但又对自己已获得的其他资源保持不放。
- 不剥夺条件：指进程已获得的资源，在未使用完之前，不能被剥夺，只能在使用完时有自己释放。
- 环路等待条件：指在发生死锁时，必然存在一个进程--资源的环形链，即进程集合{P0,P1,P2,···,Pn}中的P0正在等待一个P1占用的资源;P1正在等待P2占用的资源，······，Pn正在等待已被P0占用的资源。

### 事务A持有什么锁呢？它又想拿什么样的插入意向排他锁呢？
为了方便记录，离职用W表示Wei，J表示Jay，E表示Eason

### 我们先来分析事务A中update语句的加锁情况：
```
update  account  set balance =1000 where name ='Wei';
```
####  间隙锁
- update语句会在非唯一索引的name加上左区间的间隙锁，有区间的间隙锁（因为目前表中只有name=“Wei”的一条记录，所以没有中间的间隙锁），即（E,W）和（W，+&）
- 为什么存在间隙锁？因为这是RR的数据库隔离级别，用来解决幻读问题用的。

####  记录锁
- 因为name是索引，所以该update语句肯定会加上W的记录锁

####  Next-key锁
- Next-key锁=记录锁+间隙锁，所以该update语句就有了（E,W]的Next-key锁

### 综上所述，事务A执行完Update更新语句，会持有锁：
- Next-key Lock：（E,W]
- Cap Lock: (W,+&)

### 再来分析一下事务A中insert语句加锁的情况
```
insert into account values(null,'Jay',100);
```

### 间隙锁
- 因为Jay(J在E和W之间)，所以需要请求加（E,W）的间隙锁

### 插入意向锁（Insert Intention）
- 插入意向锁是插入一行记录操作之前设置的一种间隙锁，这个锁释放了一种插入方式的信号，即事务A需要插入意向锁（E,W）

因此，事务的update语句和insert语句执行完，它是持有了（E,w]的next-Key锁，（W,+&）的gap锁，想拿到(E,W)的插入意向排他锁。

### 事务B拥有什么间隙锁? 他为什么也要拿插入意向锁？
同理，先分析事务B的Update语句：
```
update  account  set balance =1000 where name ='Eason';
```
####  间隙锁
- update 语句会在非唯一索引的name加上左区间的间隙锁，有区间的间隙锁（因为目前表中只有name=“Eason”一套记录，所以没有中间的间隙锁），即（-&，E）和（E,W）

####  记录锁
- 因为name是索引，所以该update语句肯定会加上E的记录锁

####  Next-key锁
- Next-key Lock=记录锁+间隙锁，所以该Update语句就有了（-&，E]的Next-Key锁。

### 综上所述，事务B执行完update更行语句，会持有锁：
- Next-key Lock：（-&，E】
- Cap Lock：（E,W）

### 继续分析事务B执行insert语句加锁情况
```
insert into account  values(null,'Yan',100);
```
####  间隙锁
- 因为Yan(Y在W之后)，所以需要请求加(W,+&)的间隙锁

####  插入意向锁
- 插入意向锁是插入一行记录操作之前设置的一种间隙锁，这个锁释放了一种插入方式的信号，即事务B需要插入意向锁(W,+&)

所以，事务B的update语句和Insert语句执行完，它是持有了(-&,E)的Next-key锁，(E,W)的cap锁，想拿到(W,+&)的间隙锁，即插入意向排他锁。

### 死锁真相还原
- 事务A执行完Update Wei的语句：持有(E,W]的Next-key Lock，（W,+&）的GapLock，插入成功。
- 事务B执行完Update Eason语句，持有(-&,E]的Next-key Lock，（E,W）的GapLock，插入成功。
- 事务A执行Insert Jay的语句是，因为需要(E,W)的插入意向锁，但是(E,W)在事务B怀中，所以他陷入心塞。
- 事务B执行完Insert Yan的语句是，因为需要(W,+&)的插入意向锁，但是(W,+&)在事务A 的怀里，所以他也陷入心塞。
- 事务A持有(W,+&)的GapLock，在等待(E,W)的插入意向锁，事务B持有（E,W）的Gap锁，在等待（W,+&）的插入意向锁，所以形成了死锁的闭环。
- 事务A,B形成了死锁的闭环后，因为InnoDB的底层机制，他会让其中一个事务让出资源，另外的事务执行成功，这就是为什么事务B插入成功了，但是事务A显示DeadLock found。

# 总接
最后，遇到死锁问题，我们应该怎么分析呢？
- 模拟死锁场景
- show engine innodb status；查看死锁日志
- 找出死锁SQL
- SQL加锁分析
- 分析死锁日志
- 熟悉锁模式兼容矩阵，InnoDB存储引擎中锁的兼容矩阵。

# 3.  日常工作中mysql数据库的优化

可以从这几个维度考虑：
- 加索引
- 避免返回不必要的数据
- 适当分批量进行
- 优化SQL结构
- 分库分表
- 读写分离

##  书写高质量SQL的30条建议：
### 1.  查询SQL尽量不要使用Select *，而是select 具体字段。
理由：
- 只去需要的字段，节省资源，减少网络开销。
- select * 进行查询时，很可能就不会使用到覆盖索引了，就会造成回表查询。

### 2.  如果指导查询结果只有一条或者只要最大/最小一条记录，建议用limit 1
假设现在有employee员工表，要找出一个name叫Jay的人：
select id，name from employee where name='jay' limit 1;
理由：
- 加上limit 1后，只要找到了对应的一条记录，就不会继续乡下扫描了，效率将会大大提高。
- 当然，入股oname是唯一索引的话，是不必要加上limit 1了，因为limit的存在主要是为了防止全表扫描，从而提高性能，如果一个语句本身可以预知不用全表扫描，有没有limit，性能的差别并不大。

### 3.  营尽量避免在where子句中使用or来连接条件
理由：
- 使用or 可能回事索引失效，从而全表扫描。
```
对于Or+ 没有索引的age这种情况，假设他走了userId的索引，但是走到age查询条件时，他还得全表扫描，也就是需要三部过程：全表扫描+索引扫描+合并 如果他一开始就走全表扫描，直接一遍扫描就完事。mysql是有优化器的，处于效率与成本考虑，遇到or条件，索引可能失效，看起来也合情合理。
```

### 4. 优化limit分页
我们日常做分页需求是，一遍会用limit实现，但是当编译亮特别打的时候，查询效率就变得底下。
反例：select id，name，age from employee limit 10000，10
正例：
```
//方案一 ：返回上次查询的最大记录(偏移量)
select id，name from employee where id>10000 limit 10.

//方案二：order by + 索引
select id，name from employee order by id  limit 10000，10

//方案三：在业务允许的情况下限制页数：
```
理由：
- 当偏移量最大的时候，查询效率就会越低，因为MySQL并非是调过偏移量直接去取后面的数据，而是先把偏移量+要取的条数，然后再把前面偏移量这一段的数据抛弃掉在返回的。
- 如果使用优化方案一，返回上次最大查询记录（偏移量），这样可以调过偏移量，效率提升不少。
- 方案二使用order by+索引，也是可以提高查询效率的。
- 方案三的话，建议跟业务讨论，有没有必要查询这么后的分页拉。因为绝大多数用户都不会往后翻太多页。

### 5.  优化你的like语句
日常开发中，如果用到模糊关键字查询，很容易想到like，但是like很可能让你的索引失效。
反例：select userId，name from user where userId like '%123';
正例：select userId，name from user where userId like '123%';

### 6.  使用where条件限定要查询的数据，避免返回多余的行。
假设业务场景是这样：查询某个用户是否是会员。曾经看过老的实现代码是这样。。。
反例：List<Long> userIds = sqlMap.queryList("select userId from user where isVip=1");
boolean isVip = userIds.contains(userId);
正例：Long userId = sqlMap.queryObject("select userId from user where userId='userId' and isVip='1' ")
boolean isVip = userId！=null;

理由：
- 需要什么数据，就去查询什么数据，避免反馈不必要的数据，节省开销。

### 7.  尽量避免在索引列上使用mysql的内置函数。
业务需求：查询最近七天内登录过的用户（假设loginTime加了索引）
反例：select userId,loginTime from loginuser where Date_ADD(loginTime,Interval 7 DAY) >=now();
正例：explain  select userId,loginTime from loginuser where  loginTime >= Date_ADD(NOW(),INTERVAL - 7 DAY);

### 8.  营尽量避免在where子句中对字段进行表达式操作，这将导致系统放弃使用索引而进行全表扫描。
反例：select * from user where age-1 =10；
正例：select * from user where age =11；

### 9.  Inner Join、left Join、right join，优先使用Inner Join，入股欧式left join，左边表结果尽量小。
- Inner Join 内连接：在两张表进行连接查询时，只保留两张表中完全匹配的结果集。
- Left join 在两张表进行连接查询时，会返回左表所有的行，及时在右表中没有匹配的记录。
- Right Join 在两张表进行连接查询时，会返回右表所有的行，及时在坐标中没有匹配的记录。

都满足SQL需求的前提下，推荐优先使用Inner Join，如果要使用Left join，左边表数据结果尽量小，如果有条件的尽量放到左边处理。
反例：select * from tab1 t1 left join tab2 t2  on t1.size = t2.size where t1.id>2;
正例：select * from (select * from tab1 where id >2) t1 left join tab2 t2 on t1.size = t2.size;

理由：
- 如果Inner Join是等值连接，或许返回的行数比较小，所以性能相对会好一点。
- 同理，使用了左连接，左边表数据结果尽量小，条件尽量放到左边处理，一位着返回的行数可能比较少。

### 10. 应尽量避免在where子句中使用!= 或者<>操作符，否则将引擎放弃使用索引而进行全表扫描。
反例：select age,name  from user where age <>18;
正例：
//可以考虑分开两条sql写
select age,name  from user where age <18;
select age,name  from user where age >18;

### 11. 使用联合所以时，之一索引列的顺序，一般遵循最左匹配原则。

### 12. 对查询进行优化，应考虑在where及Order by涉及的列上建立索引，尽量避免全表扫描。

### 13. 如果插入数据过多，考虑普亮插入。
反例：
for(User u :list){
 INSERT into user(name,age) values(#name#,#age#)   
}
正例：
```
//一次500批量插入，分批进行
insert into user(name,age) values
<foreach collection="list" item="item" index="index" separator=",">
    (#{item.name},#{item.age})
</foreach>
```
理由：
- 批量插入性能好，更加节省时间。

### 14. 在适当的时候，使用覆盖索引。

### 15. 慎用distinct关键字
distinct 关键字一遍用来过滤重复记录，以返回不重复的记录。在查询一个字段或者很少的情况下使用时，给查询带来优化效果。但是在字段很多的时候使用，却会大大降低查询效率。
反例：SELECT DISTINCT * from  user;
正例：select DISTINCT name from user;

理由：
- 带distinct的语句CPU时间和占用时间都高于distinct的语句。因为当查询很多字段时，如果使用distinct，数据库引擎就会对数据库进行比较，过来掉重复数据，然而这个比较，过滤的过程会占用系统资源，CPU时间。

### 16. 删除荣誉和重复索引。
理由：宠物的索引需要维护，并且优化器在优化查询的时候也需要逐个的进行考虑，这会影响性能的。

### 17. 如果数据亮较大，优化你的修改/删除语句。
理由：
- 避免同时修改或删除过多数据，因为会造成cpu利用率过高，从而影响别人对数据库的访问。
- 一次性删除田铎数据，可能会有lock wait timeout exceed的错误，索引建议分批操作。

### 18. where子句中考虑使用默认值代替NULL。
反例：select * from user where age is not null;
正例：
//设置0为默认值
select * from user where age>0;

理由：
- 并不是说使用了is null 或者is not null就会不走索引了，这个跟mysql版本以及查询成本都有关。
- 如果mysql优化器发现，走索引比不走索引成本还高，肯定会放弃索引，这些条件!=,>is null,is not null 经常被任务索引失效，其实是应为一般情况下，查询的成本高，优化器自动放弃的。
- 如果吧NUll值，换成是默认值，很多时候让走索引成为可能，同时，表达意思会相对清晰一点。

### 19. 不要有超过5个以上的表连接。
- 连表越多，编译的时间和开销也就越大。
- 吧连接表插卡成较小的几个执行，可读性更高。
- 如果一定需要连接很多表才能得到数据，name意味着糟糕的设计了。

### 20. exist & in的合理利用
假设表A表示某企业的员工表，表B表示部门表，查询所有部门的所有员工，很容易有一下SQL：
select * from A where deptId in (select deptId from B);
这样写定价于：
```
先查询部门表B
select deptID from B
再有部门DeptID，查询A的员工
select * from A where A.deptId = B.DeptId
```
可以抽象成这样的一个循环：
```
   List<> resultSet ;
    for(int i=0;i<B.length;i++) {
          for(int j=0;j<A.length;j++) {
          if(A[i].id==B[j].id) {
             resultSet.add(A[i]);
             break;
          }
       }
    }

```
显然，出了使用in，我们也可以yogaexist实现一样的查询功能，如下：
select * from A where exist (select 1 from B where A.DeptId = B.deptUd);

因为exist查询的理解就是，先执行主查询，获得数据后，再放到子查询中做条件验证，跟进验证结果(true或者false)，来决定查询的数据结果是否得以保留。
那么，这样写就等价于：
```
select * from A；先从A表做循环
select * from B where A.deptId = B.deptId
```
同理，可以抽象成这样一个循环：
```
   List<> resultSet ;
    for(int i=0;i<A.length;i++) {
          for(int j=0;j<B.length;j++) {
          if(A[i].deptId==B[j].deptId) {
             resultSet.add(A[i]);
             break;
          }
       }
    }
```
数据库最费劲的就是跟程序连接释放。假设连接了两次，每次坐上百万次的数据集查询，查完就走，这样就只做了两次;相反见了了上百万次的连接，申请连接释放反反复复，这样那个系统就搜不了了。即mysql优化原则，就是小表驱动达标，小的数据集驱动大的数据集，从而让性能更优。

因此，选择外层循环小的，也就是，如果B的数据量小于A，适合使用in，如果B 的数据量大于A，即适合选择exist。

### 21. 尽量用union all 替换union
如果检索结果中不会有重复的记录，推荐union all 替换union。

理由：
- 如果使用union，不管检索结果有没有重复，都会尝试进行合并，然后输出最终结果前进行排序。如果已知检索结果没有重复记录，使用union all 代替union，这样会提高效率。

### 22. 索引不宜太多，一般5个以内。
- 索引并不是越多越好，索引索然提高了查询效率，但是也减低了插入和更新的效率。
- insert和update时可能会重新建索引，所以建索引需要慎重考虑，视具体情况来定。
- 一个表的索引数最好不要超过5个，若太多需要考虑一些索引是否没有存在的必要。

### 23. 尽量使用数字型字段，若只含数值信息的字段尽量不要设计为字符型。
理由：
- 相对于数字型字段，字符型会降低查询和了解的性能，并会增加存储开销。

### 24. 索引不适合健在有大量重复数据的字段上，如性别这类型数据库字段。
因为SQL优化器是跟进表中数据量来进行查询优化的，如果索引列有大量重复数据，MySQL查询优化器发现不走索引的成本更低，很可能就放弃索引了。

### 25. 尽量避免客户端返回过多数据量。
假设业务需求是，用户请求查看自己最近一年观看过得直播数据。
反例：
```
//一次性查询所有数据回来
select * from LivingInfo where watchId =useId and watchTime >= Date_sub(now(),Interval 1 Y)
```
正例：
```
//分页查询
select * from LivingInfo where watchId =useId and watchTime>= Date_sub(now(),Interval 1 Y) limit offset，pageSize

//如果是前端分页，可以先查询前两百条记录，因为一般用户应该也不会往下翻太多页，
select * from LivingInfo where watchId =useId and watchTime>= Date_sub(now(),Interval 1 Y) limit 200 ;
```
### 26. 当在SQL语句中连接多个表时，请使用表的别名，并把别名前缀于每一列上，这样语义更加清晰。
反例：
```
select  * from A inner
join B on A.deptId = B.deptId;
```
正例：
```
select  memeber.name,deptment.deptName from A member inner
join B deptment on member.deptId = deptment.deptId;
```

### 27. 尽可能使用varchar/nvarchar代替char/nchar
反例：
```
  `deptName` char(100) DEFAULT NULL COMMENT '部门名称'
```
正例：
```
  `deptName` varchar(100) DEFAULT NULL COMMENT '部门名称'
```
理由：
- 因为首先边长字段存储空间小，可节省存储空间
- 其次对于查询来说，在一个相对较小的字段内搜索，效率更高。

### 28. 为了提高group by语句的效率，可以在执行到该语句前，把不需要的记录过滤掉。
反例：
```
select job，avg（salary） from employee  group by job having job ='president' 
or job = 'managent'
```
正例：
```
select job，avg（salary） from employee where job ='president' 
or job = 'managent' group by job；
```

### 29. 如果字段类型是字符串，where时一定用引号括起来，否则索引失效。

### 30. 使用explain 分析SQL的计划。

# 4.  说说分库与分表的设计
分库分表方案，分库分表中间件，分库分表可能遇到的问题
### 1.  数据库瓶颈
不管是IO瓶颈，还是CPU瓶颈，最终都会导致数据库的活跃连接数增加，进而逼近甚至达到数据库可承载活跃连接数的阈值。在业务Service来看就是，可用数据库连接少甚至无连接可用。接下来就可以想想了（并发量、吞吐量、崩溃）。

1.  IO瓶颈
第一种：磁盘读IO瓶颈，热点数据太多，数据库缓存放不下，每次查询时会产生大量的IO，降低查询速度。-> 分库和垂直分表。
第二种：网络IO瓶颈，请求的数据太多，网络带宽不够。-> 分库。
2.  CPU瓶颈
第一种：SQL问题，如SQL中包含join、group by，order by,非索引字段条件查询等，正价CPU运算操作。->SQL 优化，建立合适的索引，在物业Service层进行业务计算。
第二中：单表数据量太大，查询时扫描的行太多，SQL效率降低，CPU率先出现瓶颈。-> 水平分表。

##  分库分表方案：
### 1.  水平分库：
- 概念：以字段为依据，按照一定策略(hash、range等)，讲一个库中的数据库拆分到多个库中。
- 结果：
  - 每个库的结构都一样；
  - 每个库的数据都不一样，没有交集。
  - 所有库的并集是全量数据。
- 场景：系统绝对并发量上来了，分表难以根本上解决问题，并且还没有明显的业务归属来垂直分库。
- 分析：库多了，IO和CPU的压力自然可以城北缓解。
### 2.  水平分表：
- 概念：以字段为依据，按照一定策略(hash、range等)，讲一个库中的数据库拆分到多个表中。
- 结果：
  - 每个表的结构都一样；
  - 每个表的数据都不一样，没有交集；
  - 所有表的并集是全量数据
- 场景：系统绝对并发量没有上来，只是单标的数据量太多了，影响了SQL效率，加重了CPU的负担，以至于成为瓶颈。
- 分析：表的数量少了，单词SQL的执行效率高，自然减轻了CPU的负担。
### 3.  垂直分库：
- 概念：以表为依据，按照业务归属不同，将不同的表拆分到不同的库中。
- 结果：
  - 每个库的结构不一样；
  - 每个库的数据也不一样，没有交集。
  - 所有库的并集是全量数据。
- 场景：系统的绝对并发量上来了，并且可以抽象出单独的业务模块。
- 分析：到这一步，基本上就可以服务化了，例如：随着业务发展一些公用的配置表、字典表等越来越多，这是可以将这些表拆到单独的库中，甚至可以服务化。再有，随着业务的发展孵化出了一套业务模式，这是可以将相关的表拆到单独的库中，甚至可以服务化。
### 4.  垂直分表：
- 概念：以字段为依据，按照字段的活跃性，将表中字段拆到不同的表中。
- 结果：
  - 每个表的结构不一样；
  - 每个表的数据也不一样，没有交集。
  - 所有表的并集是全量数据。
- 场景：系统绝对并发亮没有上来，表的记录并不多，但是字段多，并且热点数据和非热点数据在一起，单行数据库所需的存储空间较大。以至于数据库缓存的数据行减少，查询时回去读磁盘数据产生大量的随机读IO，产生IO瓶颈。
- 分析：可以用列表也和详情页来帮助理解。垂直分表的拆分远着是将热点数据（可能会瑞昱经常一起查询的数据）放在一起作为主表，非热点数据放在一起作为扩展表。这样更多的热点数据就能被缓存下来，进而减少了随机读IO。拆了之后，想要获得全部数据就需要关联两个表来读取数据。但记住，千万别用jion，因为join不仅会增加CPU负担并且会将两个表耦合在一起(必须在一个数据库实例上)。关联数据，应该在业务Service层做文章，分别获取主表和扩展表的数据，然后用关联字段得到全部数据。

##  常用的分库分表中间件：
- sharding-jdbc（当当）
- MyCat
- TDDL（淘宝）
- Oceanus（58同城）
- vitess（谷歌开发的数据库中间件）
- ATLas（Qihoo 360）

##  分库分表步骤
跟进容量(当前容量和增长量)评估分库或分表个数->选key(均匀)->分表规则(hash或range等)->执行(一般双写)->扩容问题（尽量减少数据的移动）

##  分库分表可能遇到的问题：
### 1. 事务问题
- 方案一：使用分布式事务
  - 有点：交由数据库管理，简单有效
  - 缺点：性能代价高，特别是shard越来越多时
- 方案二：由应用程序和数据库共同控制
  - 原理：将一个夸多个数据库的分布式事务分拆成多个仅处于单个数据库上面的小事务，并通过应用程序来总控各个小事务。
  - 有点：性能上有优势
  - 缺点：需要应用程序在事务控制上做灵活设计。如狗哦使用了spring的事务管理，改动起来会面临一定的困难。

### 2.  跨节点Join的问题
只要是进行切分，跨节点Join的问题是不可避免的。但是良好的设计和切分却可以减少此类情况。解决这一问题的普遍做法是分两次查询实现。在第一次查询的结果集中找出关联数据的ID，根据这些id发起第二次请求得到关联数据。

### 3.  跨节点的count，order by，group by以及聚合函数问题
这些是一类问题，因为他们都需要基于全部的数据集合进行计算。多数的代理都不会自动处理合并工作。解决方案：与解决跨节点join问题的类似，分别在各个节点上得到结果后在引用程序端进行合并，和join不同的是每个节点的查询可以并行执行，因此很多时候他的速递要比单一达标快很多，但是如果结果集很大，对应用程序内存的消耗是一个问题。

### 4.  数据迁移，容量规划，扩容等问题
来自淘宝综合业务平台团队，它利用对2的倍数取余具有向前兼容的特性（如对4取余的1的数对2取余也是1）来分配数据，避免了行级别的数据迁移，但是依然需要进行表级别的迁移，同时对扩容规模和分表数量都有限制。总的来说，这些方案都不是十分的理想，多多稍稍都存在一些缺点，这也从一个侧面反应出了sharding扩容的难点。

### 5.  事务
####  分布式事务
参考： [关于分布式事务、两阶段提交、一阶段提交、Best Efforts 1PC模式和事务补偿机制的研究](http://blog.csdn.net/bluishglc/article/details/7612811)

- 优点
1.  基于两阶段提交，最大限度的保证了跨数据库操作的原子性，是分布式系统下最严格的事务实现方式。
2.  实现简单，工作量小。由于多数应用服务器以及一些独立的分布式事务协调爱做了大量的冯作工作，是的项目中引入分布式事务的难度和工作量基本上可以忽略不急。

- 缺点
1.  系统水平伸缩的死敌。基于两阶段提交的分布式事务在提交事务时需要在多个节点之间进行协调，最大限度的推后了提交事务的时间点，客观上烟瘴了事务的实行时间，这会导致事务在访问共享资源时发生冲突和死锁的改了征稿，醉着数据库节点的增多，这种趋势会越来越严重，从而成为系统在数据库层面上水平伸缩的加锁，这是很多sharding系统不采用分布式事务的主要原因。

####  基于Best Effort 1PC模式的事务
参考spring-data-neo4j的实现。鉴于Best Effort 1PC模式的性能优势，一节相对简单的实现方式，他被大多数sharding框架和项目采用。

####  事务补偿（幂等值）
对于那些对性能要求很高，但对一致性要求不高的系统，往往并不渴求系统的实时一致性，只要在一个允许的时间周期内达到最终一致性即可，这使得事务补偿机制成为一种可行性的方案。事务补偿机制最初被提出实在“长事务”的处理中，但是对于分布式系统确保一致性也有很好的参考意义。笼统地讲，与事务在执行中发生错误后立即回滚的方式不同，事务补偿是一种时候检查并补救的措施，它只期望在一个容许的时间周期内得到最终一致的结果就可以了。事务出厂的实现与系统业务紧密相关，并没有一种标准的处理方式。一些常见的实现方式有：对数据进行队长检查；基于日志进行对比；定期同标准数据源进行同步，等等。

###  6. ID问题
一旦数据库诶切分到多个无理节点上，我们将不能再依赖数据库自身的主线生成机制。
一方面，某个分区数据库自生成的ID无法保证全局上的唯一性;另一方面，应用程序在插入数据之前需要鲜活的ID，一遍进行SQL路由。
一些常见的主线生成策略：
- UUID：
使用UUID作主线是最简单的方案，但是缺点也是非常明显的，由于UUID非常的长，除占用大量存储空间外，最主要的问题是在所以上，再建立索引和基于索引进行查询时都存在性能问题。
- 结合数据库维护一个Sequence表：
  此方案思路也很简单，在数据库中建立一个Sequence表，表结构类似于：
  ```
  CREATE TABLE `SEQUENCE` (  
    `table_name` varchar(18) NOT NULL,  
    `nextid` bigint(20) NOT NULL,  
    PRIMARY KEY (`table_name`)  
  ) ENGINE=InnoDB   
  ```

每当需要为某个表的新记录生成ID时就从Sequence表中却出对应表的nextid,并将nextid的值加1后更行的数据库中以备用下次使用。此方案也较简单，但是缺点同样明显：由于所有插入任何都需要访问该表，该表很容易成为系统性能的瓶颈，同时他也存在单点问题，一旦该表数据库失效，整个应用程序将无法工作。有人突出使用Master-slave进行主从同步，但这也只能解决单点问题，并不能解决读写比为1：1的访问压力问题。

- Twitter的分布式自增ID算法Snowflake
在分布式系统中，需要生成全局UID的场合还是比较多的，Twitter的snowflake解决了这种需求，实现还是很简单的，出去配置信息，核心代码就是毫秒级时间41位机器ID10位毫秒内序列12位。
```
* 10---0000000000 0000000000 0000000000 0000000000 0 --- 00000 ---00000 ---000000000000
```
在上面的字符串中，第一位为未使用（实际上也可以作为long的符号位），加下来41位位毫秒级时间，然后5位datacenter标识位，5位机器ID（并不算标识符，实际是为线程标识），然后12位该毫秒内的当前毫秒内的技术，加起来刚好64位，为一个Long型。

这样的好处是，整理上按照时间自增排序，并且整个分布式系统内不会产生ID碰撞（有datacenter和机器ID作区分），并且效率较高，经测试，snowflake每秒能够产生26万ID左右，完全满足需要。

### 7.  跨分片的排序分页
一般来讲，分页时需要按照制定字段进行排序。当排序字段就是片字段的时候，我们通过分片规则可以比较=容易的定位到制定的分片，而当排序字段非分片字段的时候，情况会变得比较复杂了，为了最终结果的准确性，我们需要在不同分片节点中将数据进行排序并返回，并将不同分片返回的结果集进行汇总和再次排序，最后再返回给用户。
取第一页数据：
表1：执行select ... order by date desc limit 0,10 从node1中取出前10条
表2：执行select ... order by date desc limit 0,10 从node2中取出前10条
合并再执行select ... order by date desc limit 0,10 返回最终的结果集合。

上面所描述的只是最简单的一种情况（取第一页数据），看起来对性能影响不大。但是如果想去除第10页的数据，情况有将变得复杂很多：
表1：执行select ... order by date desc limit 0,100 从node1中取出前100条
表2：执行select ... order by date desc limit 0,100 从node2中取出前100条
合并再执行select ... order by date desc limit 91,10 返回最终的结果集合。

为什么不能像获取第一页数据那样简单处理（排序取出前10条再合并、排序）。其实并不难理解，因为各分片节点中的数据可能是随机的，为了排序的准确性，必须吧所有分片节点的前N页数据都排好序后做合并，最后再进行整体的排序。很显然，这样的才做是比较消耗资源的，用户越往后翻页，系统西能将越差。

如何解决分库情况下的分页问题呢：
- 如果是在前台应用提供分页，则限定用户只能查看前面N页，这个限制在业务上是合理的，一般看后面的分页意义不大（如果一定要看，可以要求用户缩小范围重新查询）。
- 如果是后台批处理任务要求分批获取数据，则可以加大page size，比如每次获取5000条记录，有效减少分页数（淡然离线王文一般走备库，避免冲击主库）。
- 分库设计时，一般还有配套的大数据平台汇总所有分库的记录，有些分页查询可以考虑走大数据平台。

### 8.  分库策略
分库维度确定后，如何吧记录分到各个库里呢？
一般有两种方式：
- 跟进数值范围，比如用户ID为1-9999的记录分到第一个库，10000-20000的分到第二个库，以此类推。
- 根据数值驱魔，比如用户ID mod n，余数为0的记录放到第一个库，余数为1的放到第二个库，以此类推。

### 9.  分库数量
分库数量首先和单裤能处理的记录数有关，一般来说，MySQL单裤超过5000万条记录，Oracle单库超过1亿条记录，DB压力就很大，（当然处理能力和字段数量/访问模式/记录长度有进一步关系）

再满足上述前提下，如果分库数量少，达不到分散存储和减轻DB性能压力的目的；如果分库数量多，好处是每个库记录少，单库访问性能号，但对于跨多个库的访问，应用程序需要访问多个库，如果是并发模式，需要消耗宝贵的线程资源了如果是穿行模式，执行时间会急剧增加。
最后分库数量还直接影响硬件的投入，一般每个分库跑在单独无理机上，多一个库意味多一台设备。所以具体分多少个库，要综合评价，一般初次分库建议4-8个库。

### 10. 路由透明
分库从某种意义上来说，一位者DB schema改变了，必然影响应用，但这种改变和业务无关，所以要尽量保证分库对应用代码透明，分库逻辑尽量在数据访问层处理。当然完全做到这一点很困难，具体那些应该有DAL负责，那些有应用负责，这里有一些建议：
对于单裤访问，比如查询条件制定用户ID，则该SQL只需要访问特定库。此时应该由DAL层自动路由到特定库，当库二次分裂时，也只需要修改mod因子，应用代码不受影响。
对于简单的多库查询，DAL负责汇总各个数据库返回的记录，此时扔对上层应用透明。

### 11. 使用框架还是自主研发
目前市面上的分库分表中间件相对较多，其中基于代理方式的有MySQL Proxy和amoeba，基于hibername框架的是Hibernate shards，基于jdbc的有当当sharding-jdbc，基于mybatis而非类似maven插件式的有蘑菇街的蘑菇街Tsharding，通过重写sprint的ibatis template 类是cobar client，这些框架各有各的短板和优势，架构师可以再深入调研之后结合项目的实际情况进行选择，但是总的来说，我个人对于框架的选择是持谨慎态度的，一方面多数框架缺乏成功案例的验证，其成熟性与稳定性值得怀疑。另一方面，一些成功商业产品开源出矿建（如阿里和淘宝的一些开源项目）是否适合你的项目需要架构师深入调研分析的，当然，最终的选择一定是基于项目特点，团队情况、技术门槛和学习成本等综合因素确定的。

# 5.  InnoDB与MyISAM的区别
- InnoDB支持事务，MyISAM不支持事务
- InnoDB支持外键，MyISAM不支持外键
- InnoDB支持MVCC(多版本控制)，MyISAM不支持
- select count（*）from table时，MyISAM更快，因为它有一个变量保存了整个表的总行数，可以直接读取，InnoDB就需要全表扫描。
- InnoDB不支持全文索引，而MyISAM支持全文索引(5.7以后InnoDB也支持全文索引)
- InnoDB支持表、行级锁，而MyISAM支持表级锁。
- InnoDB表必须有主键，而MyISAM可以没有主键。
- InnoDB表需要更多的内存和存储，MyISAM可以被压缩，存储空间较小。
- InnoDB按主键大小有序插入，MyISAM记录插入顺序是，按记录插入顺序保存。
- InnoDB存储引擎提供了具有提交、回滚、崩溃恢复能力的事务安全，与MyISAM比InnoDB写的效率差一些，并且会占用更多的磁盘空间已保留数据和索引。

# 6.  数据库索引的原理，为什么要用B+数，为什么不用二叉树？
可以从几个维度看这个问题，查询是否够快，效率是否稳定，存储数据多少，以及查找磁盘次数，为什么不是二叉树，为什么不是平衡二叉树，为什么不是B树，而偏偏是B+数呢？

##  6.1 为什么不是一般二叉树？
如果二叉树特殊化为一个链表，相当于全表扫描。平衡二叉树相比于二叉树查找来说，查找效率更稳定，总体的查找速度也更快。

##  6.2 为什么不是平衡二叉树呢？
我们知道，在内存比在磁盘的数据，查询效率快的多。如果树这种数据结构作为索引，那我们每次查找数据就需要从磁盘中读取一个节点，也就是我们说的一个磁盘块，但是平衡二叉树可视化每个节点只存储一个键值和数据的，如果是B树，可以存储更多的节点数据，树的高度也会江都，因此读取此案的次数就降下来了，查询效率就快拉。

##  6.3 为什么不是B树而是B+树呢？
1.  B+ 树非叶子节点上是不存储数据的，仅存储键值，而B树节点中不仅存储键值，也会存储数据。InnoDB中也得默认大小是16KB，如果不存储数据，name就会存储更多的兼职，响应的数的阶数（节点的子节点树）就会更大，树就会更矮更胖，如此一来我们查找数据进行磁盘的IO次数优惠再次减少，数据查询效率也会更快。
2.  B+ 树索引的所有数据均存储在叶子节点，而且数据是按照顺序排列的，链表连着的，那么B+树是的范围查找、排序查找、分组查找以及去重查找变得异常简单。

# 7.  聚集索引与非聚集索引的区别
- 一个表中只能拥有一个聚集索引，而非聚集索引一个表可以存在多个。
- 聚集索引，索引中键值的逻辑顺序决定了表中响应的无理顺序；非聚集索引，索引中索引的顺序与磁盘上行的无理存储顺序不同。
- 索引是通过二叉树的数据结构来描述的，我们可以这么理解聚集索引：索引的叶节点就是数据节点。而非聚集索引的叶节点仍然是索引节点，只不过有一个指针只想对应的数据块。
- 聚集索引：无理存储按照索引排序；非聚集索引：无理存储不按照索引排序；

##  8.  limit 1000000加载很慢的话，你是怎么解决的？
方案一：如果ID是连续的，可以这样，返回上次查询的最大记录(偏移量)，再往下limit
```
select id，name from employee where id>1000000 limit 10.
```
方案二：在业务允许的情况下限制页数：
建议跟业务讨论，有没有必要查这么后的分页啦。因为绝大多数用户不会往后翻太多页。
方案三：order by + 索引（id为索引）
select id,name from employee order by id limit 1000000,10。
方案四：利用延迟关联或者子查询优化超多分页场景。（先快速定位需要获取的id段，然后再关联）
SELECT a.* FROM employee a, (select id from employee where 条件 LIMIT 1000000,10 ) b where a.id=b.id

# 9.  如何选择合适的分布式主线方案？
- 数据库自增长列或者字段。
- UUID
- redis生成ID
- twitter的snowflake算法
- 利用zookeeper生成唯一ID
- MongoDB的ObjectID

# 10. 事务隔离级别有那些？MySQL的默认隔离级别是什么？
- 读未提交（Read uncommited）
- 读已提交（Read commited）
- 可重复读（Repeat Read）
- 串行化（serializable）

mysql默认的事务隔离级别是可重复读（repeatable read）

# 11. 什么是幻读、脏读、不可重复读？
- 事务A、B交替执行，事务A被事务B干扰到了，因为事务A读取到事务B未提交的数据，这就是脏读。
- 在一个事务范围内，两个相同的查询，读取同一条记录，却反回了不同的数据，这就是不可重复读。
- 事务A查询一个范围结果集，另一个并发事务B往这个范围中插入/删除了数据，并静悄悄的提交，然后事务A再次查询相同的范围，两次读取得到的结果集不一样了，这就是幻读。

##  11.1  彻底读懂MySQL事务的四大隔离级别
### 事务
- 什么是事务？
  事务，有一个优先的数据库操作序列构成，这些操作要么全部执行，要么全部不执行，是一个不可分割的工作单位。
  事务的四大特性：
    原子性：事务作为一个整体被执行，包含在其中的对数据库的操作要么全部都执行，要么都不执行。
    一致性：指在事务开始之前和事务结束以后，数据不会被破坏。
    隔离性：多个事务并发访问时，事务之间是相互隔离的，一个事务不应该被其他事务干扰，多个并发事务之间要相互隔离。
    持久性：表示事务完成提交后，该事务对数据库所做的修改，将持久的保存在数据库之中。

- 事务并发存在的问题
  事务并发存在什么问题呢？欢聚话就是，一个事务是怎么干扰到其他事务的呢？
  
  假设现在有表：
  ```
  CREATE TABLE `account` (
  `id` int(11) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `balance` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `un_name_idx` (`name`) USING BTREE
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
  ```
  表中有数据：
  insert into account values (1,"Jay",100),(1,"Easoon",100),(1,"Lin",100)

  脏读：
  假设现在有AB两个事务：
  - 假设现在A的余额是100，事务A正在准备查询Jay的余额
  - 这时候，事务B先扣减Jay的余额，扣了10
  - 最后A读到的是扣减后的余额

  不可重复读：
  假设现在有AB两个事务：
  - 事务A先查询Jay的余额，查到的结果是100
  - 这时候事务B对Jay的账户余额进行扣减，扣去10后，提交事务
  - 事务A再去查询Jay的账户余额是发现变成了90

  幻读：
  假设现在有AB两个事务：
  - 事务A先查询账户ID 大于2的记录，得到了记录ID=2和3两条。
  - 这个时候，事务B开启，插入一条ID=4的记录，并且提交了
  - 事务A再去执行相同的查询，却得到了id=2/3/4。

  事务的四大隔离级别实践
  既然并发事务存在脏读、不可重复读、幻读等问题，InnoDB实现了哪几种事务的隔离级别应对呢？
  - 读未提交：存在脏读、幻读、不可重复读的问题
  - 读已提交：解决了脏读的问题
  - 可重复读：解决了不可重复读的问题
  - 串行化：解决幻读问题，性能最低

- MYSQL隔离级别的实现原理
  实现隔离机制的方法主要有两种：
  - 读写锁
  - 一致性快照读，即MVCC

  MySQL使用不同的锁策略/MVCC来实现四种不同的隔离级别。RR、RC的实现原理跟MVCC有关，RU和Serializable跟锁有关。

  读未提交，采取的是读不加锁原理：
  - 事务读不加锁，不阻塞其他书屋的读和写。
  - 事务写阻塞其他事务写，但不阻塞其他事务读。

  串行化：
  - 所有的SELECT语句会隐式的转化为SELECT ... for share，即加共享锁。
  - 读加共享锁，写加排他锁，读写互斥。如果有未提交的事务正在修改某些行，所有select这些行的语句都会阻塞。

  MVCC的实现原理
  MVCC，中文叫多版本并发控制，它是通过读取历史版本信息的数据，来境地并发事务冲突的，从而提高并发性能的一种机制。他的实现依赖于隐式字段、undo日志、快照读&、Read view，因此，我们先来了解这几个只是点。
  - 隐式字段：对于InnoDB存储引擎，每一行记录都有两个隐藏列DB_TRX_ID、DB_ROLL_PRT，如果表中没有主键和非NULL唯一键时，则还会有第三个隐藏的主键列DB_ROW_ID。
    - DB_TRX_ID，记录每一行最近一次修改(修改/更新)它的事务ID，大小为6字节。
    - DB_ROW_PTR,这个隐藏列就相当于一个指针，只想回滚端的undo日志，大小为7字节。
    - DB_ROW_ID:单调递增的行ID，带下为6字节。

  undo日志
  - 事务为提交的时候，修改数据的镜像（修改前的旧版本），存到undo日志里。以便事务回滚时，回复旧版本数据，侧小未提交事务数据对数据库的影响。
  - undo日志是逻辑日志。可以这样人物，当delete一条记录时，undo log中会记录一条对应的insert 记录，当update一条记录是，他记录一条对应相反的update记录。
  - 存储undo日志的地方，就是回滚段。
  多个事务并行操作一行数据时，不同事务对改行数据的修改会产生多个版本，然后通过指针（DB_ROLL_PTR）连一条Undo日志链。

  快照读&当前读
  - 快照读：
  读取的是记录数据的课件版本（有就的版本），不加锁，普通的select 语句都是快照读。
  - 当前读：<font color=red>读取的是记录数据的最新版本，显示加锁的都是当前读。</font>
  ```
  select * from account where id>2 lock in share mode;
  select * from  account where id>2 for update;
  ```
  Read View
  - read view 就是事务执行快照读时，产生的读视图。
  - 事务执行快照读时，会生成数据库系统当前的一个快照，记录当前系统中还有那些活跃的读写事务，把他们放到一个列表中。
  - read view主要是用来做可见性判断的，即判断当前事务可见那个版本的数据。

  为了下面方便讨论Read View可见性规则，先定义几个变量：
    - m_ids:当前系统中那些活跃的读写事务ID，他数据结构为一个List。
    - min_limit_id：m_ids事务列表中，最小的事务ID。
    - max_limit_id：m_ids事务列表中，最大的事务ID。
  - 如果DB_TRX_ID < min_limit_id，标明生成该版本的事务在生成readview前已经提交（因为事务ID是递增的）,所以该版本可以被当前事务访问。
  - 如果DB_TRX_ID > m_ids列表中最大的事务ID，标明生成该版本的事务在生成readview后才生成，所以该版本不可以被当前事务访问。
  - 如果min_limit_id =<  DB_TRX_ID <= max_limit_id,需要判断m_ids.contain(DB_TRX_ID),如果在，则代表readview生成适合，这个事务还在活跃，还有commit，你修改数据，当前事务也是看不见的，如果不在，则说明，这个事务在readview 生成之前就已经commit了，修改的结果，当前事务是能看见的。

  注意拉，RR跟RC隔离界别，最大的区别就是，RC每次读取数据前都生成一个readview，而RR只在第一次读取数据时生成一个readview。



# 12. 在高并发情况下，如何做到安全的修改同一行数据？
要安全的修改同一行数据，就要保证一个线程在修改时其他线程无法更行这行记录。一般有乐观锁和悲观锁两种。
- 使用悲观锁：
悲观锁思想就是：当前线程要来修改数据时，别的线程都得拒之门外，比如select ... for update~
- 使用乐观锁
乐观锁思想就是，有现成过来，先放过去修改。如果看到别的线程没有修改过，就可以修改成功，如果别的线程修改过，就修改失败或者重试。实现方式：乐观锁一般会使用版本号机制或CAS算法实现。

# 13. 数据库的乐观锁和悲观锁

# 14. SQL优化的一般步骤是什么，怎么看执行计划(explain)，如何理解其中各个子弹的含义。
- show status 命令理解各种SQL的执行频率
- 通过慢查询日志定位那些执行效率低的SQL语句
- explain分析低效SQL的执行计划（这点非常重要，日常开发中用它分析SQL，会大大降低SQL导致的线上事故）

##  14.1  优化SQL语句的一般步骤
### 1.  通过show status命令了解各种SQL的执行频率
MYSQL客户端连接成功后，通过show【session|global】 status 命令可以提供服务器状态信息，也可以在操作系统上使用mysqladmin extend-status命令获取这些消息，show status命令中间可以加入session(默认)或者global：
- session（当前连接）
- global（自数据库上层启动至今）
```
# Com_XXX 表示每个xxx语句执行的次数
mysql> show status like 'Com_%';
```
我们通常比较关心的是一下几个统计参数：
- Com_select：执行select操作的次数，一次查询只累加1.
- Com_insert：执行insert操作的次数，对于批量插入insert操作，只累加一次。
- Com_update：执行update操作的次数。
- Com_delete：执行Delete操作的次数。

上面这些参数对于所有存储引擎的表操作都会进行累计。下面这几个参数只是针对InnoDB的，累加的算法也略有不同：
- Innodb_rows_read：select 查询返回的行数。
- Innodb_rows_inserted：执行insert操作插入的行数。
- Innodb_rows_updated：执行update操作更行的行数。
- Innodb_rows_daleted：执行delete操作删除的行数。

通过以上几个参数，可以很容易的了解当前数据库的应用是以插入更新为主还是以查询操作为主，以及各种类型的sql大概执行比例是多少。对于更新操作的计数，是对执行次数的计数，不论提交还是回滚都会进行累加。
对于事务型的应用，通过Com_commit和Com_rollback可以了解事务提交和回滚的情况，对于回滚操作非常频繁的数据库，可能一位着应用编写存在问题。
此外，一下几个参数便于用户了解数据库的基本情况：
- Connections：试图连接mysql服务器的次数。
- Update：服务器工作时间。
- Slow_queries：慢查询次数。

### 2.  定义执行效率低的SQL语句
  1.  通过慢查询日志定位那些执行效率低的SQL语句，用--log-slow-queries[=file_name]选项启动时，mysql写一个包含所有执行时间操作long_query_time秒的sql语句的日志文件。

  2.  慢查询日志在查询结束以后才记录，所以在应用反应执行效率出现问题的时候慢查询日志并不能定位，可以使用show processlist 命令查看当前mysql在进行的线程，包括线程状态，是否锁表等，可以试试的查看SQL的执行情况，同时对一些锁表操作进行优化。  

### 3.  通过explain分析低效sql的执行计划

# 15. select for update 有什么含义，会锁表还是锁行还是其他。
它是加悲观锁，至于表锁还是行锁，这就要看是不是yoga了索引/主键了。

<font color=red>没用索引/主键的话就是表锁，否则就是行锁。</font>

# 16. Mysql 事务的四大特性以及实现原理
- 原子性：是使用undo log来实现的，如果事务直营过程中出错或者用户执行了rollback，系统就通过undo log日志返回事务开始的状态。
- 持久性，使用redo log来实现，只要redo log日志持久化了，当系统崩溃，即可通过redo log吧数据恢复。
- 隔离性：通过锁以及MVCC，使事务互相隔离开。
- 一致性：通过回滚、恢复。以及并发情况下的隔离性，从而实现一致性。

# 17. 如果某个表有近千万数据，CRUD比较慢，如何优化。
分库分表、索引优化。

# 18. 如果写SQL能够有效的使用到符合索引。
确保最佳左前缀原则。

# 19. mysql中in和exists的区别
```
select * from A where A.a in (select B.b from B ...)
select * from A where exist (select 1 from B where A.a = B.b)
```
in 先执行括号中的语句
exist先执行前面的语句
<font color=red>小表驱动大表，如果B的数量小于A，则使用in，反之使用exist</font>

# 20. 数据库自增主键可能遇到什么问题？
- 使用自增主键对数据库做分库分表，坑你出现诸如主键重复等问题。
- 自增主键会产生表锁，从而引发问题。
- 自增主键可能会用完问题。

# 21. MVCC熟悉嘛，他的底层原理？
MVCC多版本并发控制，它是通过读取历史版本的数据，来降低并发事务冲突，从而提高并发性能的一种机制。
MVCC需要关注的几个知识点：事务版本号，表的隐藏列，undo log，read view。

# 22. 数据库中间件了解过嘛？sharding jdbc，mycat。
不了解。

# 23. MYSQL的主从延迟，你怎么解决？
##  23.1  主从复制分了5个步骤进行
- 步骤1：主库的更新时间（update、insert、delete）被写到binlog中
- 步骤2：从库发起连接，连接到主库
- 步骤3：此时主库创建一个binlog dump thread,把binlog 的内容发送到从库。
- 步骤4：从库启动之后，创建一个IO线程，读取主库传过来的binlog内容并写入到relay log。
- 步骤5：还会创建一个SQL线程，从relaylog里面读取内容。从exec_master_log_pos位置开始执行读取到的更行事件，将更新内容写入到slave的DB。

##  23.2  主从同步延迟的原因
一个服务器开放N个连接给客户端来连接的，这样会有大并发的更新操作，但是从服务器的里面读取binlog的线程仅有一个，当某个SQL在从服务器上执行的事件稍长，或者由于某个SQL要进行锁表就会导致，主服务器的SQL大量积压，未被同步到从服务器里。这就导致了主从不一致，也就是主从延迟。

##  23.3  主从延迟的解决办法
- 主服务器要负责更新操作，对安全性的要求比从服务器的要搞，所以有些设置参数可以修改，比如sync_binlog=1,innodb_flush_log_at_trx_commit=1之类的设置等。
- 选择更好的硬件设备作为slave。
- 把一台从服务器当作备份使用，而不提供查询，那么他的负载下来了，执行relay log里面的SQL效率自然就搞了。
- 增加从服务器，目的是分散读的压力，从而降低服务器负载。

# 24. 锁一下大表查询的优化方案：
- 优化schema、sql语句+索引；
- 可以考虑加缓存，memcached，redis，或者JVM本地缓存；
- 主从复制，读写分离；
- 分库分表。

# 25. 什么是数据库连接池？为什么需要数据库连接池呢？
- 数据库连接池基本原理：在内部对象池中，维护一定数量的数据库连接，并对外暴露数据库连接的获取和返回方法。
- 应用程序和数据库建立连接的过程：
  - 通过TCP协议的三次握手和数据库服务器建立连接
  - 发送数据库用户账号密码，等待数据库验证用户身份
  - 完成身份验证后，系统可以提交SQL语句到数据库执行
  - 把连接关闭，停tcp四次挥手告别。

- 数据库连接池的好处：
  - 资源重用（连接重复）:
  - 更快的系统响应速冻
  - 新的资源分配手段
  - 统一的连接管理，避免数据库连接泄露

# 26. 一条SQL语句在MYSQL中如何执行的
查询语句：
- 先检查该语句是否有权限；
- 如果没有权限，直接返回错误信息；
- 如果有权限，在MYSQL8.0版本以前，会先查询缓存。8.0以后移除。
- 如果没有缓存，分析器进行词法分析，提前SQL语句select等关键元素。然后判读SQL语句是否有语法错误，比如关键词是否正确等等。
- 优化器进行确定执行方案。
- 进行权限校验，如果没有权限就直接返回错误信息，如果有权限就会调用数据库引擎接口，返回执行结果。

# 27. InnoDB引擎中的索引策略，了解过嘛？
- 覆盖索引
- 最佳左前缀原则
- <font color=red>索引下推</font>

索引下推优化是MYSQL5.6引入的，可以在索引遍历过程中，对索引中包含的字段先做判断，直接过滤掉不满足的巨鹿，减少回表次数。

# 28. 数据库存储日期格式是，如果考虑失去转换问题？
- datetime类型适合用来记录数据的原始创建时间，修改记录中其他的值，datetime字段的值不会改变。
- timestamp类型适合用来记录数据的最后修改时间，只要修改了记录中其他字段的值，timestamp子弹的值都会被自动更新。

# 29. 一条sql执行过长的时间，你如果优化，从哪方面入手？
- 查看是否涉及多表和子查询，优化sql结构，如去除荣誉字段，是否可拆表等
- 优化索引结构，看是否可以适当添加索引
- 数量打的表，可以考虑进行分离/分表（如交易流水等）
- 数据库主从复制，读写分离
- explain分析sql语句，查看执行计划，优化SQL
- 查看mysql执行日志，分析是否有其他方面的问题。

# 30. MYSQL数据库服务器性能分析方法命令有哪些？
- showstatus，一些值得监控的变量值：
  - Bytes_received和Bytes_sent和服务器之间来往的流量
  - Com_*服务器正在执行的命令。
  - Created_*在查询执行期限间创建的临时表和文件。
  - Handler_*存储引擎操作
  - Select_*不同类型的连接执行计划
  - Sort_*几种排序信息。

- show profiles是MYSQL用来分析当前回话SQL语句执行的资源消耗情况。

# 31. Blob和text有什么区别？
- Blob踊跃存储二进制数据，而Text用于存储大字符串。
- Blob值被视为二进制字符串（字节字符串），他们没有字符集，并且排序和比较基于列值中的字节数值。
- text值被视为非二进制字符串（字符字符串）。他们有一个字符集，并根据字符集的排序规则对值进行排序和比较。

# 32. MYSQL里面记录货币用什么字段类型比较号？
- 货币在数据库中MYSQL从用Decimal和NUMric类型表示，这两种类型被MYSQL实现为同样的类型。他们被用于保存与金钱有关的数据。
- salary DECIMAL（9，2），9（precision）代表江北用于存储值的小数位数，而2（scale）代表将被踊跃存储小数点后的位数。存储在salary列中的值的范围是从-9999999.99到9999999.99。
- DECIMAL和NUMERIC值作为字符串存储，而不是作为二进制浮点数，一遍保存那些值的小数精度。

# 33. MYSQL中有哪几种锁？
表锁、行锁、页锁

# 34. Hash索引和B+树区别是什么？你在设计索引是怎么抉择的？
- B+树可以进行范围查询，hash不行
- B+树支持联合索引的最左侧原则，Hash索引不支持
- B+树支持Order by 排序，Hash不支持
- Hash 索引在等值查询上比B+树效率更高
- B+树使用like进行模糊查询的时候，like后面（比如%）的话可以起到优化作用，hash无法进行模糊查询。

# 35. MYSQL的内连接、左连接、右连接的区别
- Inner join 内连接，在两张表进行连接查询时，只保留两张表中完全匹配的结果集
- left join 在两张表进行连接查询时，会返回左表所有的行，即使在右表中没有匹配的记录。
- right join 在两张表进行连接查询时，会返回右表所有的行，即使在左表中没有匹配的记录。

# 