CREATE TABLE IF NOT EXISTS firends (
    id INT AUTO_INCREMENT PRIMARY KEY,
    friends JSON,
    userId INT,
    CONSTRAINT fk_user2
    FOREIGN key (userId)
    REFERENCES user_login(userId)
)


CREATE TABLE IF NOT EXISTS friend_request (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_id INT,
    to_id INT,
    request_status ENUM('0', '1')
)