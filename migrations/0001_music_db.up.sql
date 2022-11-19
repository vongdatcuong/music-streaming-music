CREATE TABLE IF NOT EXISTS Genre (
	`genre_id` INT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(256) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,

    PRIMARY KEY (`genre_id`)
);

CREATE TABLE IF NOT EXISTS Song (
	`song_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(256) NOT NULL,
    `genre` INT(20) UNSIGNED NOT NULL,
    `artist` VARCHAR(256) NOT NULL,
    `duration` INT UNSIGNED NOT NULL, 
    `language` VARCHAR(10) NOT NULL,
    `rating` FLOAT DEFAULT 0,
    `resource_id` CHAR(36) NOT NULL,
    `resource_link` VARCHAR(10000) NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `status` TINYINT UNSIGNED NOT NULL,
    
    PRIMARY KEY (`song_id`),
    FOREIGN KEY (`genre`) REFERENCES Genre(`genre_id`)
);

CREATE TABLE IF NOT EXISTS Playlist (
	`playlist_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(256) NOT NULL,
    `created_by` BIGINT(20) UNSIGNED NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `status` TINYINT UNSIGNED NOT NULL,
    
    PRIMARY KEY (`playlist_id`)
);

CREATE TABLE IF NOT EXISTS Playlist_Song (
	`playlist_id` BIGINT(20) UNSIGNED NOT NULL,
    `song_id` BIGINT(20) UNSIGNED NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    
    FOREIGN KEY (`playlist_id`) REFERENCES Playlist(`playlist_id`),
    FOREIGN KEY (`song_id`) REFERENCES Song(`song_id`)
);

INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Drama', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());
INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Kpop', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());
INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Jpop', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());
INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Chinese', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());
INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Hip-hop', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());
INSERT INTO Genre (`name`, `created_at`, `updated_at`) VALUES ('Romantic', UNIX_TIMESTAMP(),  UNIX_TIMESTAMP());