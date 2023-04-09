
-- Create the table to store the CO2 readings
CREATE TABLE IF NOT EXISTS `co2` (
  `id` INTEGER PRIMARY KEY,
  `ppm` INTEGER NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the table to store the temperature, pressure and humidity readings
CREATE TABLE IF NOT EXISTS `bme280` (
  `id` INTEGER PRIMARY KEY,
  `temperature` REAL NOT NULL,
  `pressure` REAL NOT NULL,
  `humidity` REAL NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create the table to store the person-in-room movements
CREATE TABLE IF NOT EXISTS `pir` (
  `id` INTEGER PRIMARY KEY,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

