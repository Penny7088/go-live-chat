
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
    FOREIGN KEY (native_language_id) REFERENCES languages(id) ON DELETE SET NULL,
    FOREIGN KEY (learning_language_id) REFERENCES languages(id) ON DELETE SET NULL,
    FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE SET NULL
);

-- 语言表 (languages)
CREATE TABLE languages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    native_name VARCHAR(100) NOT NULL,
    iso_code VARCHAR(10) UNIQUE NOT NULL
);

-- 国家表 (countries)
CREATE TABLE countries (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    iso_code VARCHAR(10) UNIQUE NOT NULL
);

-- 国家语言关联表 (country_languages)
CREATE TABLE country_languages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    country_id BIGINT NOT NULL,
    language_id BIGINT NOT NULL,
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
    INDEX (sender_id),
    INDEX (receiver_id),
    INDEX (timestamp),
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 离线消息表 (offline_messages)
CREATE TABLE offline_messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,
    message TEXT NOT NULL,
    message_type ENUM('text', 'image', 'audio', 'video') NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_delivered BOOLEAN DEFAULT FALSE,
    INDEX (receiver_id),
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
    provider ENUM('google', 'facebook', 'apple') NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
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
    device_token VARCHAR(255) NOT NULL,
    device_type ENUM('iOS', 'Android', 'Web') NOT NULL,
    last_active TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


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
INSERT INTO countries (name, iso_code) VALUES
                                           ('Australia', 'AU'),
                                           ('Canada', 'CA'),
                                           ('Norway', 'NO'),
                                           ('Netherlands', 'NL'),
                                           ('Iceland', 'IS'),
                                           ('Indonesia', 'ID'),
                                           ('Philippines', 'PH'),
                                           ('Singapore', 'SG'),
                                           ('United States', 'US'),
                                           ('Mexico', 'MX'),
                                           ('France', 'FR'),
                                           ('Germany', 'DE'),
                                           ('China', 'CN'),
                                           ('India', 'IN'),
                                           ('Brazil', 'BR'),
                                           ('Russia', 'RU'),
                                           ('Japan', 'JP'),
                                           ('United Kingdom', 'GB'),
                                           ('Italy', 'IT'),
                                           ('South Korea', 'KR'),
                                           ('Spain', 'ES'),
                                           ('Thailand', 'TH'),
                                           ('Vietnam', 'VN'),
                                           ('Ukraine', 'UA'),
                                           ('Saudi Arabia', 'SA'),
                                           ('Egypt', 'EG'),
                                           ('Argentina', 'AR'),
                                           ('Switzerland', 'CH'),
                                           ('Austria', 'AT'),
                                           ('Malaysia', 'MY');

INSERT INTO country_languages (country_id, language_id) VALUES
-- Australia - English
((SELECT id FROM countries WHERE name = 'Australia'), (SELECT id FROM languages WHERE iso_code = 'en')),
-- Canada - English, French
((SELECT id FROM countries WHERE name = 'Canada'), (SELECT id FROM languages WHERE iso_code = 'en')),
-- Norway - Norwegian (Note: Norwegian language may need to be added)
((SELECT id FROM countries WHERE name = 'Norway'), (SELECT id FROM languages WHERE iso_code = 'no')),
-- Netherlands - Dutch (Note: Dutch language may need to be added)
((SELECT id FROM countries WHERE name = 'Netherlands'), (SELECT id FROM languages WHERE iso_code = 'nl')),
-- Iceland - Icelandic (Note: Icelandic language may need to be added)
((SELECT id FROM countries WHERE name = 'Iceland'), (SELECT id FROM languages WHERE iso_code = 'is')),
-- Indonesia - Indonesian
((SELECT id FROM countries WHERE name = 'Indonesia'), (SELECT id FROM languages WHERE iso_code = 'id')),
-- Philippines - Filipino, English
((SELECT id FROM countries WHERE name = 'Philippines'), (SELECT id FROM languages WHERE iso_code = 'fil')),
-- Singapore - English, Chinese (Simplified), Malay, Tamil
((SELECT id FROM countries WHERE name = 'Singapore'), (SELECT id FROM languages WHERE iso_code = 'en'))
-- Continue for other countries and languages as needed...
