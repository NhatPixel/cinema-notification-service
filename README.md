//SQL tạo bảng

CREATE TABLE notifications (
    id CHAR(36) NOT NULL,
    user_id VARCHAR(50) NOT NULL,
    title VARCHAR(255),
    content TEXT,
    PRIMARY KEY (id),
    INDEX idx_user_id (user_id)
);
