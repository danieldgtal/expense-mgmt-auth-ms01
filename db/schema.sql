-- Enable UUID extension for generating unique IDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create USER table
CREATE TABLE IF NOT EXISTS "users" (
  user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  salt VARCHAR(50) NOT NULL,
  is_email_verified BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  reset_password_token VARCHAR(255),
  reset_password_expires TIMESTAMP WITH TIME ZONE,
  currency CHAR(3) DEFAULT 'CAD' NOT NULL CHECK (currency IN ('CAD', 'USD'))
);

-- Create AUTH_TOKEN table
CREATE TABLE IF NOT EXISTS auth_tokens (
  token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
  token VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create REFRESH_TOKEN table
-- CREATE TABLE IF NOT EXISTS refresh_tokens (
--     refresh_token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
--     token VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     expires_at TIMESTAMP WITH TIME ZONE NOT NULL
-- );

-- -- Create CATEGORY table
-- CREATE TABLE IF NOT EXISTS categories (
--   category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--   user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
--   name VARCHAR(50) NOT NULL,
--   description TEXT,
--   color_code VARCHAR(7),
--   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--   deleted_at TIMESTAMP WITH TIME ZONE
-- );
CREATE TABLE IF NOT EXISTS categories (
  category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID REFERENCES "users"(user_id) ON DELETE CASCADE,  -- Allow null for default categories
  name VARCHAR(50) NOT NULL,
  description TEXT,
  color_code VARCHAR(7),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);


-- Create BUDGET table
CREATE TABLE IF NOT EXISTS budgets (
  budget_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(category_id) ON DELETE SET NULL,
  amount DECIMAL(10, 2) CHECK (amount >= 0) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create RECEIPT table
CREATE TABLE IF NOT EXISTS receipts (
  receipt_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  image_url VARCHAR(255) NOT NULL,
  ocr_data TEXT,
  scanned_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create EXPENSE table
CREATE TABLE IF NOT EXISTS expenses (
  expense_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(category_id) ON DELETE SET NULL,
  amount DECIMAL(10, 2) CHECK (amount >= 0) NOT NULL,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  description TEXT,
  receipt_id UUID REFERENCES receipts(receipt_id) ON DELETE SET NULL,
  is_recurring BOOLEAN DEFAULT FALSE,
  recurring_interval_days INT CHECK (recurring_interval_days >= 0),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create USER_SETTINGS table with default currency set to CAD
CREATE TABLE IF NOT EXISTS user_settings (
  settings_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
  notifications_enabled BOOLEAN DEFAULT TRUE,
  language_preference VARCHAR(10),
  theme_preference VARCHAR(20)
);

-- Create NOTIFICATION table
CREATE TABLE IF NOT EXISTS notifications (
  notification_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES "users"(user_id) ON DELETE CASCADE,
  type VARCHAR(50) NOT NULL,
  message TEXT NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create AUDIT_LOG table to track changes in critical tables
CREATE TABLE IF NOT EXISTS audit_logs (
  log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID,
  table_name VARCHAR(50) NOT NULL,
  action VARCHAR(50) NOT NULL,
  old_data JSONB,
  new_data JSONB,
  change_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES "users"(user_id) ON DELETE SET NULL
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_expense_user_id ON expenses(user_id);
CREATE INDEX IF NOT EXISTS idx_expense_date ON expenses(date);
CREATE INDEX IF NOT EXISTS idx_budget_user_id ON budgets(user_id);
CREATE INDEX IF NOT EXISTS idx_category_user_id ON categories(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_user_id ON notifications(user_id);
