-- 000001_initial_schema.down.sql
DROP TABLE IF EXISTS attachments;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_skills;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS conversation_participants;
DROP TABLE IF EXISTS conversations;
DROP TABLE IF EXISTS order_bids;
DROP TABLE IF EXISTS order_tags;
DROP TABLE IF EXISTS portfolio_files;
DROP TABLE IF EXISTS portfolios;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS user_types;

DROP TYPE IF EXISTS notification_type;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS bid_status;
DROP TYPE IF EXISTS order_status;
