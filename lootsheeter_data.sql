-- --------------------------------------------------------
-- Host:                         network.morpheusxaut.net
-- Server version:               5.5.40-0ubuntu0.14.04.1 - (Ubuntu)
-- Server OS:                    debian-linux-gnu
-- HeidiSQL Version:             9.1.0.4882
-- --------------------------------------------------------

USE `lootsheeter`;

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
-- Dumping data for table lootsheeter.fleetroles: ~35 rows (approximately)
/*!40000 ALTER TABLE `fleetroles` DISABLE KEYS */;
INSERT IGNORE INTO `fleetroles` (`id`, `ship`, `fleet_role`) VALUES
	(1, 'vexor navy issue', 32),
	(2, 'tengu', 32),
	(3, 'loki', 32),
	(4, 'legion', 32),
	(5, 'proteus', 32),
	(6, 'huggin', 32),
	(7, 'rapier', 32),
	(8, 'moa', 32),
	(9, 'ferox', 32),
	(10, 'ishtar', 32),
	(11, 'drake', 32),
	(12, 'dominix', 32),
	(13, 'pilgrim', 32),
	(14, 'orthrus', 32),
	(15, 'gila', 32),
	(16, 'bellicose', 64),
	(17, 'scimitar', 16),
	(18, 'scythe', 16),
	(19, 'noctis', 8),
	(20, 'catalyst', 8),
	(21, 'cormorant', 8),
	(22, 'heron', 4),
	(23, 'magnate', 4),
	(24, 'anathema', 4),
	(25, 'probe', 4),
	(26, 'buzzard', 4),
	(27, 'imicus', 4),
	(28, 'helios', 4),
	(29, 'cheetah', 4),
	(30, 'purifier', 4),
	(31, 'manticore', 4),
	(32, 'nemesis', 4),
	(33, 'hound', 4),
	(34, 'thrasher', 8),
	(35, 'vexor', 32);
/*!40000 ALTER TABLE `fleetroles` ENABLE KEYS */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
