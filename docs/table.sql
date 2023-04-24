/*
SQLyog Community v13.1.6 (64 bit)
MySQL - 5.7.33 : Database - teach
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;
/*!40101 SET SQL_MODE = ''*/;
/*!40014 SET @OLD_UNIQUE_CHECKS = @@UNIQUE_CHECKS, UNIQUE_CHECKS = 0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES = @@SQL_NOTES, SQL_NOTES = 0 */;

CREATE DATABASE /*!32312 IF NOT EXISTS */`teach` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;

USE `teach`;

/*Table structure for table `t_class` */

DROP TABLE IF EXISTS `t_class`;

CREATE TABLE `t_class`
(
    `id`       bigint(20)   NOT NULL AUTO_INCREMENT,
    `class_id` varchar(64)  NOT NULL COMMENT '班级编号',
    `name`     varchar(128) NOT NULL COMMENT '班级姓名',
    `college`  varchar(128) NOT NULL COMMENT '学院名称',
    `major`    varchar(128) NOT NULL COMMENT '专业',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_teacher_id` (`class_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_exam` */

DROP TABLE IF EXISTS `t_exam`;

CREATE TABLE `t_exam`
(
    `id`                bigint(20)   NOT NULL AUTO_INCREMENT,
    `exam_id`           varchar(64)  NOT NULL COMMENT '试卷编号',
    `exam_name`         varchar(128) NOT NULL COMMENT '试卷名称',
    `questions`         text         NOT NULL COMMENT '题⽬, 序列化字符串',
    `every_score`       varchar(1024)  DEFAULT NULL COMMENT '试卷题目对应的分值',
    `comment`           varchar(255) NOT NULL COMMENT '备注',
    `create_teacher_id` varchar(64)  NOT NULL COMMENT '创建教师编号',
    `update_time`       datetime     NOT NULL COMMENT '修改时间',
    `create_time`       datetime     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_knowledge_point` */

DROP TABLE IF EXISTS `t_knowledge_point`;

CREATE TABLE `t_knowledge_point`
(
    `id`            bigint(20)   NOT NULL AUTO_INCREMENT,
    `knp_id`        varchar(64)  NOT NULL COMMENT '知识点编号',
    `name`          varchar(128) NOT NULL COMMENT '知识点名称',
    `parent_knp_id` varchar(64)  NOT NULL COMMENT '⽗知识点编号',
    `level`         tinyint(4)   NOT NULL COMMENT '困难程度,1: 容易，2: 中等, 3:困难',
    `content`       text         NOT NULL COMMENT '知识点内容',
    `create_user`   varchar(64)  NOT NULL COMMENT '录⼊者',
    `update_time`   datetime     NOT NULL COMMENT '修改时间',
    `create_time`   datetime     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_knowledge_point_question` */

DROP TABLE IF EXISTS `t_knowledge_point_question`;

CREATE TABLE `t_knowledge_point_question`
(
    `id`          bigint(20)  NOT NULL AUTO_INCREMENT,
    `knp_id`      varchar(64) NOT NULL COMMENT '知识点编号编号',
    `question_id` varchar(64) NOT NULL COMMENT '题⽬编号',
    `update_time` datetime    NOT NULL COMMENT '修改时间',
    `create_time` datetime    NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_question` */

DROP TABLE IF EXISTS `t_question`;

CREATE TABLE `t_question`
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
    `question_id` varchar(64)  NOT NULL COMMENT '题⽬编号',
    `name`        varchar(128) NOT NULL COMMENT '题⽬名称',
    `knp_id`      varchar(64)  NOT NULL COMMENT '上级知识点编号',
    `level`       tinyint(4)   NOT NULL COMMENT '困难程度,1: 容易，2: 中等, 3:困难',
    `type`        tinyint(4)   NOT NULL COMMENT '题⽬类型, 1:选择题, 2:填空题, 3:问答题',
    `content`     text         NOT NULL COMMENT '题⽬内容',
    `answer`      text         NOT NULL COMMENT '题⽬答案',
    `create_user` varchar(64)  NOT NULL COMMENT '录⼊者',
    `update_time` datetime     NOT NULL COMMENT '修改时间',
    `create_time` datetime     NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_student` */

DROP TABLE IF EXISTS `t_student`;

CREATE TABLE `t_student`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT,
    `student_id`   varchar(64)  NOT NULL COMMENT '学⽣编号',
    `name`         varchar(128) NOT NULL COMMENT '学⽣姓名',
    `password`     varchar(255) NOT NULL COMMENT '密码, 加盐哈希',
    `college`      varchar(128) NOT NULL COMMENT '学院名称',
    `major`        varchar(128) NOT NULL COMMENT '专业',
    `class_id`     bigint(20)   NOT NULL COMMENT '班级编号',
    `phone_number` varchar(64)  NOT NULL COMMENT '电话号码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_student_id` (`student_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_student_exam` */

DROP TABLE IF EXISTS `t_student_exam`;

CREATE TABLE `t_student_exam`
(
    `id`             bigint(20)  NOT NULL AUTO_INCREMENT,
    `exam_id`        varchar(64) NOT NULL COMMENT '试卷编号',
    `student_id`     varchar(64) NOT NULL COMMENT '学⽣编号',
    `answers`        text COMMENT '答案, 序列化字符串',
    `comment`        varchar(255)         DEFAULT NULL COMMENT '备注',
    `start_time`     datetime    NOT NULL COMMENT '考试开始时间',
    `end_time`       datetime    NOT NULL COMMENT '考试结束时间',
    `detailed_score` varchar(1024)        DEFAULT NULL COMMENT '试卷题目详细分数',
    `total_score`    float(10, 1)         DEFAULT NULL COMMENT '总分数',
    `problem_status` varchar(1024)        DEFAULT NULL COMMENT '题目状态',
    `status`         tinyint(1)  NOT NULL DEFAULT '-1' COMMENT '状态: -1未提交 0未批改完 1未批改  2已批改完成',
    `times`          tinyint(2)           DEFAULT '1' COMMENT '做的次数',
    `update_time`    datetime    NOT NULL COMMENT '修改时间',
    `create_time`    datetime    NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_student_question` */

DROP TABLE IF EXISTS `t_student_question`;

CREATE TABLE `t_student_question`
(
    `id`          bigint(20)  NOT NULL AUTO_INCREMENT,
    `question_id` varchar(64) NOT NULL COMMENT '题⽬编号',
    `student_id`  varchar(64) NOT NULL COMMENT '学⽣编号',
    `homework_id` varchar(64) NOT NULL COMMEN                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                  满分',
    `score`       float       NOT NULL COMMENT '分值，得分',
    `answer`      text        NOT NULL COMMENT '答案, 序列化字符串',
    `update_time` datetime    NOT NULL COMMENT '修改时间',
    `create_time` datetime    NOT NULL COMMENT '创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_teacher` */

DROP TABLE IF EXISTS `t_teacher`;


CREATE TABLE `t_teacher`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT,
    `teacher_id`   varchar(64)  NOT NULL COMMENT '教师编号',
    `name`         varchar(128) NOT NULL COMMENT '教师姓名',
    `password`     varchar(255) NOT NULL COMMENT '密码, 加盐哈希',
    `college`      varchar(128) NOT NULL COMMENT '学院名称',
    `major`        varchar(128) NOT NULL COMMENT '专业',
    `phone_number` varchar(64)  NOT NULL COMMENT '电话号码',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_teacher_id` (`teacher_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*Table structure for table `t_teacher_class` */

DROP TABLE IF EXISTS `t_teacher_class`;


CREATE TABLE `t_teacher_class`
(
    `id`         bigint(20)  NOT NULL AUTO_INCREMENT,
    `class_id`   varchar(64) NOT NULL COMMENT '班级编号',
    `teacher_id` varchar(64) NOT NULL COMMENT '班级编号',
    `is_valid`   tinyint(4)  NOT NULL COMMENT '是否合法, 0: 不合法, 1:合法',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_teacher_class` (`teacher_id`, `class_id`),
    KEY `idx_class` (`class_id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

/*!40101 SET SQL_MODE = @OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS = @OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES = @OLD_SQL_NOTES */;
