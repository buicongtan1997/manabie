CREATE DATABASE manabie;
CREATE DATABASE manabie_test;

CREATE USER 'manabie'@'%' identified BY 'manabie';
FLUSH PRIVILEGES;

GRANT ALL PRIVILEGES ON manabie.* TO 'manabie'@'%';
GRANT ALL PRIVILEGES ON manabie_test.* TO 'manabie'@'%';
