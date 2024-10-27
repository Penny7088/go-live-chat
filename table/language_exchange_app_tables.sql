-- 用户表 (users)
CREATE TABLE users (
                       id BIGINT AUTO_INCREMENT PRIMARY KEY,
                       email VARCHAR(255) UNIQUE,
                       username VARCHAR(100) NOT NULL,
                       password_hash VARCHAR(255),
                       profile_picture VARCHAR(255),
                       native_language_id BIGINT,
                       learning_language_id BIGINT,
                       language_level VARCHAR(50),
                       age INT,
                       gender ENUM('male', 'female', 'other'),
                       interests TEXT,
                       country_id BIGINT,
                       registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       last_login TIMESTAMP,
                       status ENUM('active', 'inactive', 'banned') DEFAULT 'active',
                       email_verified BOOLEAN DEFAULT FALSE,
                       verification_token VARCHAR(255),
                       token_expiration TIMESTAMP,
                       created_at datetime     null,
                       updated_at datetime     null,
                       deleted_at datetime     null,
                       FOREIGN KEY (native_language_id) REFERENCES languages(id) ON DELETE SET NULL,
                       FOREIGN KEY (learning_language_id) REFERENCES languages(id) ON DELETE SET NULL,
                       FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE SET NULL
);

-- 群组聊天表 (group_chats)
CREATE TABLE group_chats (
                             id BIGINT AUTO_INCREMENT PRIMARY KEY,
                             group_name VARCHAR(255) NOT NULL,
                             created_by BIGINT NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- 群组成员表 (group_members)
CREATE TABLE group_members (
                               id BIGINT AUTO_INCREMENT PRIMARY KEY,
                               group_id BIGINT NOT NULL,
                               user_id BIGINT NOT NULL,
                               joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               is_admin BOOLEAN DEFAULT FALSE,
                               FOREIGN KEY (group_id) REFERENCES group_chats(id) ON DELETE CASCADE,
                               FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 群组消息表 (group_messages)
CREATE TABLE group_messages (
                                id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                group_id BIGINT NOT NULL,
                                sender_id BIGINT NOT NULL,
                                message TEXT NOT NULL,
                                message_type ENUM('text', 'image', 'audio', 'video') NOT NULL,
                                timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                FOREIGN KEY (group_id) REFERENCES group_chats(id) ON DELETE CASCADE,
                                FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 邮箱验证表 (email_verifications)
CREATE TABLE email_verifications (
                                     id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                     user_id BIGINT NOT NULL,
                                     verification_token VARCHAR(255),
                                     token_expiration TIMESTAMP,
                                     is_verified BOOLEAN DEFAULT FALSE,
                                     sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 语言表 (languages)
CREATE TABLE languages (
                           id BIGINT AUTO_INCREMENT PRIMARY KEY,
                           name VARCHAR(100) NOT NULL,
                           native_name VARCHAR(100) NOT NULL,
                           iso_code VARCHAR(10) UNIQUE NOT NULL,
                           created_at datetime     null,
                           updated_at datetime     null,
                           deleted_at datetime     null
);

-- 国家表 (countries)
CREATE TABLE countries (
                           id BIGINT AUTO_INCREMENT PRIMARY KEY,
                           name VARCHAR(100) NOT NULL,
                           iso_code VARCHAR(10) UNIQUE NOT NULL,
                           created_at datetime     null,
                           updated_at datetime     null,
                           deleted_at datetime     null
);

-- 国家语言关联表 (country_languages)
CREATE TABLE country_languages (
                                   id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                   country_id BIGINT NOT NULL,
                                   language_id BIGINT NOT NULL,
                                   created_at datetime     null,
                                   updated_at datetime     null,
                                   deleted_at datetime     null,
                                   FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE,
                                   FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE
);

-- 消息表 (messages)
CREATE TABLE messages (
                          id BIGINT AUTO_INCREMENT PRIMARY KEY,
                          sender_id BIGINT NOT NULL,
                          receiver_id BIGINT NOT NULL,
                          message TEXT NOT NULL,
                          message_type ENUM('text', 'image', 'audio', 'video') NOT NULL,
                          timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          is_delivered BOOLEAN DEFAULT FALSE,
                          is_read BOOLEAN DEFAULT FALSE,
                          INDEX (sender_id),
                          INDEX (receiver_id),
                          INDEX (timestamp),
                          FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
                          FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户聊天记录备份表 (user_chat_backup)
CREATE TABLE user_chat_backup (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  user_id BIGINT NOT NULL,
                                  chat_partner_id BIGINT NOT NULL,
                                  message TEXT NOT NULL,
                                  message_type ENUM('text', 'image', 'audio', 'video') NOT NULL,
                                  timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  is_restored BOOLEAN DEFAULT FALSE,
                                  INDEX (user_id),
                                  INDEX (chat_partner_id),
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                  FOREIGN KEY (chat_partner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 词汇表和笔记 (vocabulary_notes)
CREATE TABLE vocabulary_notes (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  user_id BIGINT NOT NULL,
                                  word VARCHAR(255) NOT NULL,
                                  note TEXT,
                                  language_id BIGINT NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                                  FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE
);

-- 语音房间表 (voice_chat_rooms)
CREATE TABLE voice_chat_rooms (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  room_name VARCHAR(255) NOT NULL,
                                  language_id BIGINT NOT NULL,
                                  created_by BIGINT NOT NULL,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE CASCADE,
                                  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- 第三方登录表 (third_party_auth)
CREATE TABLE third_party_auth (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  user_id BIGINT NOT NULL,
                                  provider VARCHAR(50) NOT NULL,
                                  provider_user_id VARCHAR(255) NOT NULL,
                                  created_at datetime     null,
                                  updated_at datetime     null,
                                  deleted_at datetime     null,
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 虚拟货币系统表 (virtual_currency)
CREATE TABLE virtual_currency (
                                  id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                  user_id BIGINT NOT NULL,
                                  balance INT DEFAULT 0,
                                  last_transaction_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 消息队列表 (message_queue)
CREATE TABLE message_queue (
                               id BIGINT AUTO_INCREMENT PRIMARY KEY,
                               sender_id BIGINT NOT NULL,
                               receiver_id BIGINT NOT NULL,
                               message TEXT NOT NULL,
                               message_type ENUM('text', 'image', 'audio', 'video') NOT NULL,
                               timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               is_sent BOOLEAN DEFAULT FALSE,
                               retry_count INT DEFAULT 0,
                               INDEX (receiver_id),
                               FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
                               FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 推送通知日志表 (push_notification_log)
CREATE TABLE push_notification_log (
                                       id BIGINT AUTO_INCREMENT PRIMARY KEY,
                                       user_id BIGINT NOT NULL,
                                       notification_type ENUM('message', 'friend_request', 'system') NOT NULL,
                                       content TEXT NOT NULL,
                                       timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       is_clicked BOOLEAN DEFAULT FALSE,
                                       INDEX (user_id),
                                       FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 用户设备表 (user_devices)
CREATE TABLE user_devices (
                              id BIGINT AUTO_INCREMENT PRIMARY KEY,
                              user_id BIGINT NOT NULL,
                              device_token VARCHAR(255) NOT NULL UNIQUE,
                              device_type ENUM('iOS', 'Android', 'Web') NOT NULL,
                              ip_address VARCHAR(255)  NOT NULL,
                              created_at datetime     null,
                              updated_at datetime     null,
                              deleted_at datetime     null,
                              last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              INDEX (user_id),
                              FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 预插入数据 (languages, countries, country_languages)
INSERT INTO languages (name, native_name, iso_code) VALUES
                                                        ('Arabic', 'العربية', 'ar'),
                                                        ('German', 'Deutsch', 'de'),
                                                        ('English', 'English', 'en'),
                                                        ('Portuguese', 'Português', 'pt'),
                                                        ('French', 'Français', 'fr'),
                                                        ('Chinese Simplified', '简体中文', 'zh-Hans'),
                                                        ('Chinese Traditional', '繁體中文', 'zh-Hant'),
                                                        ('Italian', 'Italiano', 'it'),
                                                        ('Japanese', '日本語', 'ja'),
                                                        ('Korean', '한국어', 'ko'),
                                                        ('Spanish', 'Español', 'es'),
                                                        ('Russian', 'Русский', 'ru'),
                                                        ('Thai', 'ไทย', 'th'),
                                                        ('Vietnamese', 'Tiếng Việt', 'vi'),
                                                        ('Ukrainian', 'Українська', 'uk');

INSERT INTO Countries (name, iso_code) VALUES
                                                       ('United States', 'USA'),
                                                       ('Canada', 'CAN'),
                                                       ('United Kingdom', 'GBR'),
                                                       ('Australia', 'AUS'),
                                                       ('New Zealand', 'NZL'),
                                                       ('Germany', 'DEU'),
                                                       ('France', 'FRA'),
                                                       ('Italy', 'ITA'),
                                                       ('Spain', 'ESP'),
                                                       ('Netherlands', 'NLD'),
                                                       ('Belgium', 'BEL'),
                                                       ('Sweden', 'SWE'),
                                                       ('Norway', 'NOR'),
                                                       ('Denmark', 'DNK'),
                                                       ('Finland', 'FIN'),
                                                       ('Ireland', 'IRL'),
                                                       ('Switzerland', 'CHE'),
                                                       ('Austria', 'AUT'),
                                                       ('Portugal', 'PRT'),
                                                       ('Greece', 'GRC'),
                                                       ('Poland', 'POL'),
                                                       ('Czech Republic', 'CZE'),
                                                       ('Hungary', 'HUN'),
                                                       ('Romania', 'ROU'),
                                                       ('Russia', 'RUS'),
                                                       ('China', 'CHN'),
                                                       ('Japan', 'JPN'),
                                                       ('South Korea', 'KOR'),
                                                       ('India', 'IND'),
                                                       ('Brazil', 'BRA'),
                                                       ('Mexico', 'MEX'),
                                                       ('Argentina', 'ARG'),
                                                       ('Chile', 'CHL'),
                                                       ('Colombia', 'COL'),
                                                       ('South Africa', 'ZAF'),
                                                       ('Egypt', 'EGY'),
                                                       ('Nigeria', 'NGA'),
                                                       ('Kenya', 'KEN'),
                                                       ('Turkey', 'TUR'),
                                                       ('Saudi Arabia', 'SAU'),
                                                       ('United Arab Emirates', 'ARE'),
                                                       ('Israel', 'ISR'),
                                                       ('Indonesia', 'IDN'),
                                                       ('Thailand', 'THA'),
                                                       ('Vietnam', 'VNM'),
                                                       ('Philippines', 'PHL'),
                                                       ('Malaysia', 'MYS'),
                                                       ('Singapore', 'SGP'),
                                                       ('Bangladesh', 'BGD');


INSERT INTO country_languages (country_id, language_id) VALUES
                                                            (1, 3), (2, 3), (3, 3), (4, 3), (5, 3),
                                                            (6, 3), (7, 3), (8, 3), (9, 3), (10, 11),
                                                            (11, 5), (12, 2), (13, 6), (13, 7), (14, 3),
                                                            (14, 6), (15, 4), (16, 12), (17, 9), (18, 3),
                                                            (19, 8), (20, 10), (21, 11), (22, 13), (23, 15), (24, 14);
