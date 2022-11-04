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
    `resource_id` BIGINT(20) UNSIGNED NOT NULL,
    `resource_link` VARCHAR(10000),
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `status` INT UNSIGNED NOT NULL,
    
    PRIMARY KEY (`song_id`),
    FOREIGN KEY (`genre`) REFERENCES Genre(`genre_id`)
);

CREATE TABLE IF NOT EXISTS Playlist (
	`playlist_id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
    `name` VARCHAR(256) NOT NULL,
    `created_by` BIGINT(20) UNSIGNED NOT NULL,
    `created_at` BIGINT(20) NOT NULL,
    `updated_at` BIGINT(20) NOT NULL,
    `status` INT UNSIGNED NOT NULL,
    
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