-- Enable UUID extension for generating unique IDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create USER table
CREATE TABLE IF NOT EXISTS "user" (
    user_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    salt VARCHAR(50) NOT NULL,
    is_email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    reset_password_token VARCHAR(255),
    reset_password_expires TIMESTAMP WITH TIME ZONE
);

-- Create AUTH_TOKEN table
CREATE TABLE IF NOT EXISTS auth_token (
    token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create REFRESH_TOKEN table
CREATE TABLE IF NOT EXISTS refresh_token (
    refresh_token_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Create CATEGORY table
CREATE TABLE IF NOT EXISTS category (
    category_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    color_code VARCHAR(7)
);

-- Create CURRENCY table for currency management
CREATE TABLE IF NOT EXISTS currency (
    currency_code CHAR(3) PRIMARY KEY,
    currency_name VARCHAR(50) NOT NULL
);

-- Insert default currencies (USD and CAD) if they do not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM currency WHERE currency_code = 'USD') THEN
        INSERT INTO currency (currency_code, currency_name) VALUES ('USD', 'United States Dollar');
    END IF;

    IF NOT EXISTS (SELECT 1 FROM currency WHERE currency_code = 'CAD') THEN
        INSERT INTO currency (currency_code, currency_name) VALUES ('CAD', 'Canadian Dollar');
    END IF;
END $$;

-- Create BUDGET table
CREATE TABLE IF NOT EXISTS budget (
    budget_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    category_id UUID REFERENCES category(category_id) ON DELETE SET NULL,
    amount DECIMAL(10, 2) CHECK (amount >= 0) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    currency_code CHAR(3) NOT NULL REFERENCES currency(currency_code) ON DELETE RESTRICT
);

-- Create RECEIPT table
CREATE TABLE IF NOT EXISTS receipt (
    receipt_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_url VARCHAR(255) NOT NULL,
    ocr_data TEXT,
    scanned_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create EXPENSE table
CREATE TABLE IF NOT EXISTS expense (
    expense_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    category_id UUID REFERENCES category(category_id) ON DELETE SET NULL,
    amount DECIMAL(10, 2) CHECK (amount >= 0) NOT NULL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    description TEXT,
    receipt_id UUID REFERENCES receipt(receipt_id) ON DELETE SET NULL,
    currency_code CHAR(3) NOT NULL REFERENCES currency(currency_code) ON DELETE RESTRICT,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurring_interval_days INTEGER CHECK (recurring_interval_days >= 0)
);

-- Create USER_SETTINGS table with default currency set to CAD
CREATE TABLE IF NOT EXISTS user_settings (
    settings_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    default_currency CHAR(3) NOT NULL DEFAULT 'CAD' REFERENCES currency(currency_code) ON DELETE RESTRICT,
    notifications_enabled BOOLEAN DEFAULT TRUE,
    language_preference VARCHAR(10),
    theme_preference VARCHAR(20)
);

-- Create NOTIFICATION table
CREATE TABLE IF NOT EXISTS notification (
    notification_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create AUDIT_LOG table to track changes in critical tables
CREATE TABLE IF NOT EXISTS audit_log (
    log_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    table_name VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    old_data JSONB,
    new_data JSONB,
    change_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user"(user_id) ON DELETE SET NULL
);

-- Create indexes for frequently queried fields
CREATE INDEX IF NOT EXISTS idx_expense_user_id ON expense(user_id);
CREATE INDEX IF NOT EXISTS idx_expense_date ON expense(date);
CREATE INDEX IF NOT EXISTS idx_budget_user_id ON budget(user_id);
CREATE INDEX IF NOT EXISTS idx_category_user_id ON category(user_id);
CREATE INDEX IF NOT EXISTS idx_notification_user_id ON notification(user_id);
