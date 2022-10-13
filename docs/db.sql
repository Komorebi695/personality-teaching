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

CREATE TABLE `t_class` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `class_id` varchar(64) NOT NULL COMMENT '班级编号',
  `name` varchar(128) NOT NULL COMMENT '班级姓名',
  `college` varchar(128) NOT NULL COMMENT '学院名称',
  `major` varchar(128) NOT NULL COMMENT '专业',
  `is_delete` int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

/*Data for the table `t_class` */

insert  into `t_class`(`id`,`class_id`,`name`,`college`,`major`,`is_delete`) values (1,'1','计科01','计科院','计科',0);
insert  into `t_class`(`id`,`class_id`,`name`,`college`,`major`,`is_delete`) values (2,'31','代角长事十','elit dolor dolore adipisicing enim','eiusmod laboris reprehenderit',0);
insert  into `t_class`(`id`,`class_id`,`name`,`college`,`major`,`is_delete`) values (3,'82','三活更重音员少','eu do','quis',1);

/*Table structure for table `t_student` */

CREATE TABLE `t_student` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `student_id` varchar(64) NOT NULL COMMENT '学生编号',
  `name` varchar(128) NOT NULL COMMENT '学生姓名',
  `password` varchar(255) NOT NULL COMMENT '密码，加盐哈希',
  `college` varchar(128) NOT NULL COMMENT '学院名称',
  `major` varchar(128) NOT NULL COMMENT '专业',
  `class_id` varchar(64) NOT NULL COMMENT '班级编号',
  `phone_number` varchar(64) NOT NULL COMMENT '电话号码',
  `is_delete` int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_student_id` (`student_id`),
  UNIQUE KEY `uniq_phone_number` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

/*Data for the table `t_student` */

insert  into `t_student`(`id`,`student_id`,`name`,`password`,`college`,`major`,`class_id`,`phone_number`,`is_delete`) values (1,'202021091195','我去二','123456','计科院','计科','1','111111111',0);
insert  into `t_student`(`id`,`student_id`,`name`,`password`,`college`,`major`,`class_id`,`phone_number`,`is_delete`) values (2,'22222','号号步市认立','minim cupidatat consectetur quis','minim ut officia occaecat','consequat id proident','1','85',0);
insert  into `t_student`(`id`,`student_id`,`name`,`password`,`college`,`major`,`class_id`,`phone_number`,`is_delete`) values (3,'77','切记条今始候','ex consectetur Duis','eu aliquip','sunt quis','1','73',0);
insert  into `t_student`(`id`,`student_id`,`name`,`password`,`college`,`major`,`class_id`,`phone_number`,`is_delete`) values (4,'22','员领则半验领今','ut dolore id','exercitation ut','eiusmod minim','1','86',1);

/*Table structure for table `t_teacher` */

CREATE TABLE `t_teacher` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `teacher_id` varchar(64) NOT NULL COMMENT '教师编号',
  `name` varchar(128) NOT NULL COMMENT '教师姓名',
  `password` varchar(255) NOT NULL COMMENT '密码，加盐哈希',
  `college` varchar(128) NOT NULL COMMENT '学院名称',
  `major` varchar(128) NOT NULL COMMENT '专业',
  `phone_number` varchar(64) NOT NULL COMMENT '电话号码',
  `is_delete` int(8) NOT NULL COMMENT '是否已删除；0：否；1：是',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_teacher_id` (`teacher_id`),
  UNIQUE KEY `uniq_phone_number` (`phone_number`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

/*Data for the table `t_teacher` */

insert  into `t_teacher`(`id`,`teacher_id`,`name`,`password`,`college`,`major`,`phone_number`,`is_delete`) values (1,'74','义进保期','eu ut','计算机科学学院','计科','2222222222',0);

/*Table structure for table `t_teacher_class` */

CREATE TABLE `t_teacher_class` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `class_id` varchar(64) NOT NULL COMMENT '班级编号',
  `teacher_id` varchar(64) NOT NULL COMMENT '教师编号',
  `is_valid` tinyint(4) NOT NULL COMMENT '是否合法；0：不合法；1：合法',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_teacher_class` (`teacher_id`,`class_id`),
  KEY `idx_class` (`class_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `t_teacher_class` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
