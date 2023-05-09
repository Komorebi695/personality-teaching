/*
SQLyog Ultimate v12.08 (64 bit)
MySQL - 5.7.19 : Database - school_teach
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*Table structure for table `t_class` */

CREATE TABLE `t_class`
(
    `id`        bigint(20) NOT NULL AUTO_INCREMENT,
    `class_id`  varchar(64)  NOT NULL COMMENT '班级编号',
    `name`      varchar(128) NOT NULL COMMENT '班级姓名',
    `college`   varchar(128) NOT NULL COMMENT '学院名称',
    `major`     varchar(128) NOT NULL COMMENT '专业',
    `is_delete` int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

/*Table structure for table `t_student` */

CREATE TABLE `t_student`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT,
    `student_id`   varchar(64)  NOT NULL COMMENT '学生编号',
    `name`         varchar(128) NOT NULL COMMENT '学生姓名',
    `password`     varchar(255) NOT NULL COMMENT '密码，加盐哈希',
    `college`      varchar(128) NOT NULL COMMENT '学院名称',
    `major`        varchar(128) NOT NULL COMMENT '专业',
    `class_id`     varchar(64)  NOT NULL COMMENT '班级编号',
    `phone_number` varchar(64)  NOT NULL COMMENT '电话号码',
    `is_delete`    int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_student_id` (`student_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

/*Table structure for table `t_teacher` */

CREATE TABLE `t_teacher`
(
    `id`           bigint(20) NOT NULL AUTO_INCREMENT,
    `teacher_id`   varchar(64)  NOT NULL COMMENT '教师编号',
    `name`         varchar(128) NOT NULL COMMENT '教师姓名',
    `password`     varchar(255) NOT NULL COMMENT '密码，加盐哈希',
    `college`      varchar(128) NOT NULL COMMENT '学院名称',
    `major`        varchar(128) NOT NULL COMMENT '专业',
    `phone_number` varchar(64)  NOT NULL COMMENT '电话号码',
    `is_delete`    int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_teacher_id` (`teacher_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Table structure for table `t_teacher_class` */

CREATE TABLE `t_teacher_class`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT,
    `class_id`   varchar(64) NOT NULL COMMENT '班级编号',
    `teacher_id` varchar(64) NOT NULL COMMENT '教师编号',
    `is_valid`   tinyint(4) NOT NULL COMMENT '是否合法；0：不合法；1：合法',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_teacher_class` (`teacher_id`,`class_id`),
    KEY          `idx_class` (`class_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `t_teacher_class` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

