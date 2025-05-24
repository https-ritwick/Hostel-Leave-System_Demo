-- Drop existing tables if any
DROP TABLE IF EXISTS user_login;
DROP TABLE IF EXISTS leave_requests;
DROP TABLE IF EXISTS user_details;

-- User Details Table
CREATE TABLE user_details (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    contact_number VARCHAR(15),
    gender ENUM('Male', 'Female') NOT NULL,
    room_number VARCHAR(10),
    address TEXT,
    role ENUM('Student', 'Warden', 'UWD') DEFAULT 'Student',
    validflag BOOLEAN DEFAULT 1
);

-- Login Credentials Table
CREATE TABLE user_login (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    userpassword VARCHAR(255) NOT NULL,
    validflag BOOLEAN DEFAULT 1,
    user_details_ref_id INT,
    FOREIGN KEY (user_details_ref_id) REFERENCES user_details(id) ON DELETE CASCADE
);

-- Leave Requests Table
CREATE TABLE leave_requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT NOT NULL,
    reason TEXT NOT NULL,
    from_date DATE NOT NULL,
    to_date DATE NOT NULL,
    status ENUM('Pending', 'Approved', 'Rejected') DEFAULT 'Pending',
    FOREIGN KEY (student_id) REFERENCES user_details(id) ON DELETE CASCADE
);
