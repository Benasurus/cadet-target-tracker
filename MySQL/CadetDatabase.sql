-- MySQL dump 10.13  Distrib 5.7.12, for Win64 (x86_64)
--
-- Host: localhost    Database: cadettracker
-- ------------------------------------------------------
-- Server version	5.7.14-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `sectors`
--

DROP TABLE IF EXISTS `sectors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sectors` (
  `sectorID` varchar(64) NOT NULL,
  `sectorLabel` varchar(64) DEFAULT NULL,
  `sectorDescription` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`sectorID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sectors`
--

LOCK TABLES `sectors` WRITE;
/*!40000 ALTER TABLE `sectors` DISABLE KEYS */;
INSERT INTO `sectors` VALUES ('sec1-set0','Drill','Drill'),('sec1-set1','First Class Standard','First Class Standard'),('sec1-set2','Local Parade','Local Parade'),('sec1-set3','Drill Competition','Drill Competition'),('sec1-set4','Drill and Ceremonial','Drill and Ceremonial'),('sec10-set0','Leadership','Leadership'),('sec10-set1','SMEAC','SMEA'),('sec10-set2','JNCO Course','JNCO Course'),('sec10-set3','SCNO Course','SNCO Course'),('sec10-set4','National Leadership Course','National Leadership Course'),('sec11-set0','DofE','DofE'),('sec11-set1','Bronze','Bronze'),('sec11-set2','Silver','Silver'),('sec11-set3','Gold','Gold'),('sec11-set4','DofE Supervisor','DofE Supervisor'),('sec12-set0','Community Engagement','Community Engagement'),('sec12-set1','Charity Collection','Charity Collection'),('sec12-set2','Squadron Recruitment','Squadron Recruitment'),('sec12-set3','Wing Parade','Wing Parade'),('sec12-set4','Lord Lieutenant\'s Cadet','Lord Lieutenant\'s Cadet'),('sec13-set0','Shooting','Shooting'),('sec13-set1','No.8 WHT','No.8 WHT'),('sec13-set2','L98A2 WHT','L98A2 WHT'),('sec13-set3','Marksmenship Badges','Marksmenship Badges'),('sec13-set4','Bisley/Cadet 100','Bisley/Cadet 100'),('sec14-set0','Music','Music'),('sec14-set1','Band Practice','Band Practice'),('sec14-set2','Band Badge','Band Badge'),('sec14-set3','Music Camp','Music Camp'),('sec14-set4','Band Competition','Band Competition'),('sec15-set0','Camps','Camps'),('sec15-set1','Squadron','Squadron'),('sec15-set2','Wing Camp','Wing Camp'),('sec15-set3','Annual Camp','Annual Camp'),('sec15-set4','International Camp','International Camp'),('sec2-set0','Radio','Radio'),('sec2-set1','Basic Radio Operator','Basic Radio Operator'),('sec2-set2','Radio Operator','Radio Operator'),('sec2-set3','Communicator','Communicator'),('sec2-set4','Communications Specialist','Communications Specialist'),('sec3-set0','Flying','Flying'),('sec3-set1','Airmenship','Airmenship'),('sec3-set2','Flight Simulator','Flight Simulator'),('sec3-set3','AEF Flight','AEF Flight'),('sec3-set4','Flying Scholarship','Flying Scholarship'),('sec4-set0','Gliding','Gliding'),('sec4-set1','GIC 1','GIC 1'),('sec4-set2','GIC 2','GIC 2'),('sec4-set3','Gliding Scholarship','Gliding Scholarship'),('sec4-set4','Gliding Instructor','Gliding Instructor'),('sec5-set0','Fieldcraft','Fieldcraft'),('sec5-set1','ACP 16 Lessons','ACP 16 Lessons'),('sec5-set2','Fieldcraft Exercise','Fieldcraft Exercise'),('sec5-set3','Fieldcraft weekend','Fieldcraft weekend'),('sec5-set4','Fieldcraft Instructor','Fieldcraft Instructor'),('sec6-set0','Classifications','Classifications'),('sec6-set1','First Class','First Class'),('sec6-set2','Leading','Leading'),('sec6-set3','Senior/Master','Senior/Master'),('sec6-set4','Instructor Cadet','Instructor Cadet'),('sec7-set0','Sports','Sports'),('sec7-set1','Squadron Sports','Squadron Sports'),('sec7-set2','Inter Squadron Sport','Inter Squadron Sport'),('sec7-set3','Inter Wing Sport','Inter Wing Sport'),('sec7-set4','Corps Sport','Corps Sport'),('sec8-set0','AT','AT'),('sec8-set1','IET','IET'),('sec8-set2','Day Activity','Day Activity'),('sec8-set3','AT Weekend','AT Weekend'),('sec8-set4','National Camp','National Camp'),('sec9-set0','First Aid','First Aid'),('sec9-set1','Heartstart','Heartstart'),('sec9-set2','YFA','YFA'),('sec9-set3','First Aid Competition','First Aid Competition'),('sec9-set4','AFA','AFA');
/*!40000 ALTER TABLE `sectors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `transactions` (
  `transactionID` int(11) NOT NULL AUTO_INCREMENT,
  `userName` varchar(64) DEFAULT NULL,
  `timeStamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `section1` int(1) DEFAULT '0',
  `section2` int(1) DEFAULT '0',
  `section3` int(1) DEFAULT '0',
  `section4` int(1) DEFAULT '0',
  `section5` int(1) DEFAULT '0',
  `section6` int(1) DEFAULT '0',
  `section7` int(1) DEFAULT '0',
  `section8` int(1) DEFAULT '0',
  `section9` int(1) DEFAULT '0',
  `section10` int(1) DEFAULT '0',
  `section11` int(1) DEFAULT '0',
  `section12` int(1) DEFAULT '0',
  `section13` int(1) DEFAULT '0',
  `section14` int(1) DEFAULT '0',
  `section15` int(1) DEFAULT '0',
  `cpi` float DEFAULT NULL,
  PRIMARY KEY (`transactionID`),
  KEY `userName` (`userName`),
  CONSTRAINT `transactions_ibfk_1` FOREIGN KEY (`userName`) REFERENCES `userdata` (`userName`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
INSERT INTO `transactions` VALUES (1,'jeff','2016-11-16 00:00:01',0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,27),(2,'dave','2016-11-18 00:00:01',0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,38),(3,'dave','2016-12-01 00:00:01',0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,40),(4,'dave','2016-12-11 00:00:01',0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,43),(5,'bob','2016-12-11 00:00:05',0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,38);
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userdata`
--

DROP TABLE IF EXISTS `userdata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `userdata` (
  `userName` varchar(64) NOT NULL,
  `firstName` varchar(64) DEFAULT NULL,
  `surnameName` varchar(64) DEFAULT NULL,
  `dateOfBirth` date DEFAULT NULL,
  `dateOfEnrollment` date DEFAULT NULL,
  `sex` char(1) DEFAULT NULL,
  `flight` char(1) DEFAULT NULL,
  PRIMARY KEY (`userName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userdata`
--

LOCK TABLES `userdata` WRITE;
/*!40000 ALTER TABLE `userdata` DISABLE KEYS */;
INSERT INTO `userdata` VALUES ('bob','bob','',NULL,NULL,'M','A'),('dave','dave','',NULL,NULL,'M','B'),('jeff','jeff','',NULL,NULL,'M','C');
/*!40000 ALTER TABLE `userdata` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2017-01-14 17:18:12
