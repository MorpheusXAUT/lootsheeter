-- --------------------------------------------------------
-- Host:                         network.morpheusxaut.net
-- Server version:               5.5.40-0ubuntu0.14.04.1 - (Ubuntu)
-- Server OS:                    debian-linux-gnu
-- HeidiSQL Version:             9.1.0.4880
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- Dumping database structure for lootsheeter
CREATE DATABASE IF NOT EXISTS `lootsheeter` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `lootsheeter`;


-- Dumping structure for table lootsheeter.corporations
CREATE TABLE IF NOT EXISTS `corporations` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `corporation_id` bigint(20) NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `ticker` varchar(10) COLLATE utf8_unicode_ci DEFAULT NULL,
  `active` enum('Y','N') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Y',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `corp_id` (`corporation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Data exporting was unselected.


-- Dumping structure for table lootsheeter.fleetmembers
CREATE TABLE IF NOT EXISTS `fleetmembers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `fleet_id` bigint(20) NOT NULL,
  `player_id` bigint(20) NOT NULL,
  `role` int(10) NOT NULL DEFAULT '0',
  `site_modifier` int(10) NOT NULL DEFAULT '0',
  `payment_modifier` double NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `fleet_id_player_id` (`fleet_id`,`player_id`),
  KEY `fk_fleetmembers_player` (`player_id`),
  KEY `fk_fleetmembers_fleet` (`fleet_id`),
  CONSTRAINT `fk_fleetmembers_fleet` FOREIGN KEY (`fleet_id`) REFERENCES `fleets` (`id`),
  CONSTRAINT `fk_fleetmembers_player` FOREIGN KEY (`player_id`) REFERENCES `players` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Data exporting was unselected.


-- Dumping structure for table lootsheeter.fleets
CREATE TABLE IF NOT EXISTS `fleets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `system` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `system_nickname` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `profit` double NOT NULL DEFAULT '0',
  `losses` double NOT NULL DEFAULT '0',
  `sites_finished` int(10) NOT NULL DEFAULT '0',
  `start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NULL DEFAULT NULL,
  `active` enum('Y','N') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Y',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Data exporting was unselected.


-- Dumping structure for table lootsheeter.players
CREATE TABLE IF NOT EXISTS `players` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `player_id` bigint(20) NOT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `corporation_id` bigint(20) NOT NULL,
  `access` int(10) NOT NULL DEFAULT '0',
  `active` enum('Y','N') COLLATE utf8_unicode_ci NOT NULL DEFAULT 'Y',
  PRIMARY KEY (`id`),
  UNIQUE KEY `player_id` (`player_id`),
  UNIQUE KEY `name` (`name`),
  KEY `fk_players_corporation` (`corporation_id`),
  CONSTRAINT `fk_players_corporation` FOREIGN KEY (`corporation_id`) REFERENCES `corporations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- Data exporting was unselected.
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
