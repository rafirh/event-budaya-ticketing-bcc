-- =====================================================
-- DATABASE: EVENT BUDAYA MALANG
-- PostgreSQL UUID VERSION
-- =====================================================

-- =====================================================
-- EXTENSION
-- =====================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- =====================================================
-- USERS
-- =====================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(150) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    phone VARCHAR(20),

    role VARCHAR(20) NOT NULL
        CHECK (role IN ('user','promotor','admin')),

    profile_photo TEXT,
    gender VARCHAR(10)
        CHECK (gender IN ('male','female','other')),

    email_verified_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- EMAIL VERIFICATION TOKENS
-- =====================================================
CREATE TABLE email_verification_tokens (
    id SERIAL PRIMARY KEY,

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    token_hash VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_email_verification_user ON email_verification_tokens(user_id);
CREATE INDEX idx_email_verification_expires ON email_verification_tokens(expires_at);


-- =====================================================
-- EVENT CATEGORIES
-- =====================================================
CREATE TABLE event_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    icon TEXT
);


-- =====================================================
-- EVENTS
-- =====================================================
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    promoter_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    category_id UUID
        REFERENCES event_categories(id),

    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) UNIQUE,

    summary VARCHAR(255),
    description TEXT,

    venue VARCHAR(200),
    address TEXT,
    google_maps_url TEXT,

    start_date TIMESTAMP,
    end_date TIMESTAMP,
    registration_deadline TIMESTAMP,

    is_paid BOOLEAN DEFAULT TRUE,
    price NUMERIC(12,2) DEFAULT 0,
    quota INT NOT NULL DEFAULT 0 CHECK (quota >= 0),
    sold INT DEFAULT 0 CHECK (sold >= 0),

    banner_url TEXT,

    status VARCHAR(20) DEFAULT 'draft'
        CHECK (status IN ('draft','published','finished')),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- ORDERS
-- =====================================================
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL
        REFERENCES users(id),

    event_id UUID NOT NULL
        REFERENCES events(id) ON DELETE CASCADE,

    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(12,2) NOT NULL,
    service_fee NUMERIC(12,2) DEFAULT 0,
    total_price NUMERIC(12,2) DEFAULT 0,

    status VARCHAR(20) DEFAULT 'pending'
        CHECK (status IN ('pending','paid','cancelled')),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- PAYMENTS
-- =====================================================
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    payment_type VARCHAR(50) NOT NULL DEFAULT 'ORDER'
        CHECK (payment_type IN ('ORDER', 'EVENT_POSTING_FEE')),

    order_id UUID
        REFERENCES orders(id) ON DELETE CASCADE,

    event_id UUID
        REFERENCES events(id) ON DELETE CASCADE,

    payment_method VARCHAR(50),
    payment_gateway VARCHAR(50),
    payment_url TEXT,

    amount NUMERIC(12,2),

    status VARCHAR(20)
        CHECK (status IN ('waiting','success','failed')),

    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_event ON payments(event_id);
CREATE INDEX idx_payments_type ON payments(payment_type);


-- =====================================================
-- TICKETS (E-TICKET)
-- =====================================================
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    order_id UUID NOT NULL
        REFERENCES orders(id) ON DELETE CASCADE,

    ticket_code VARCHAR(100) UNIQUE NOT NULL,
    qr_code TEXT,

    holder_name VARCHAR(150) NOT NULL,
    identity_type VARCHAR(50) NOT NULL,
    identity_number VARCHAR(100) NOT NULL,
    holder_phone VARCHAR(20) NOT NULL,
    holder_email VARCHAR(150) NOT NULL,
    notes TEXT NOT NULL,

    is_used BOOLEAN DEFAULT FALSE,
    used_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW()
);


-- =====================================================
-- EVENT COMMENTS (DISCUSSION)
-- =====================================================
CREATE TABLE event_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    event_id UUID NOT NULL
        REFERENCES events(id) ON DELETE CASCADE,

    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    parent_id UUID
        REFERENCES event_comments(id) ON DELETE CASCADE,

    comment TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- PROMOTER WALLETS
-- =====================================================
CREATE TABLE promoter_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    promoter_id UUID UNIQUE NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    balance NUMERIC(14,2) DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW()
);


-- =====================================================
-- WALLET TRANSACTIONS
-- =====================================================
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    wallet_id UUID NOT NULL
        REFERENCES promoter_wallets(id) ON DELETE CASCADE,

    type VARCHAR(50) NOT NULL
        CHECK (type IN ('TICKET_COMMISSION','WITHDRAW')),

    direction VARCHAR(10) NOT NULL
        CHECK (direction IN ('IN','OUT')),

    amount NUMERIC(12,2) NOT NULL,

    reference_id UUID,
    description TEXT,

    created_at TIMESTAMP DEFAULT NOW()
);


-- =====================================================
-- INDEXES (PERFORMANCE)
-- =====================================================
CREATE INDEX idx_users_role ON users(role);

CREATE INDEX idx_events_promoter ON events(promoter_id);
CREATE INDEX idx_events_category ON events(category_id);
CREATE INDEX idx_events_status ON events(status);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_event ON orders(event_id);
CREATE INDEX idx_orders_status ON orders(status);

CREATE INDEX idx_tickets_order ON tickets(order_id);
CREATE INDEX idx_tickets_identity_number ON tickets(identity_number);

CREATE INDEX idx_comments_event ON event_comments(event_id);

CREATE INDEX idx_wallet_promoter ON promoter_wallets(promoter_id);


-- =====================================================
-- FEE SETTINGS (SERVICE FEE & EVENT POSTING FEE)
-- =====================================================
CREATE TABLE fee_settings (
    id SERIAL PRIMARY KEY,

    fee_type VARCHAR(50) NOT NULL
        CHECK (fee_type IN ('SERVICE_FEE', 'EVENT_POSTING_FEE')),

    calculation_type VARCHAR(20) NOT NULL
        CHECK (calculation_type IN ('percentage', 'fixed')),

    amount NUMERIC(12, 2) NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- ADMIN WALLET
-- =====================================================
CREATE TABLE admin_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    balance NUMERIC(14, 2) DEFAULT 0,

    total_revenue NUMERIC(14, 2) DEFAULT 0,
    total_withdrawn NUMERIC(14, 2) DEFAULT 0,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);


-- =====================================================
-- PROMOTER TRANSACTION HISTORY
-- =====================================================
CREATE TABLE promoter_transaction_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    promoter_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    transaction_type VARCHAR(50) NOT NULL
        CHECK (transaction_type IN ('TICKET_COMMISSION', 'WITHDRAW', 'POSTING_FEE', 'REFUND', 'OTHER')),

    direction VARCHAR(10) NOT NULL
        CHECK (direction IN ('IN', 'OUT')),

    amount NUMERIC(12, 2) NOT NULL,
    balance_before NUMERIC(14, 2),
    balance_after NUMERIC(14, 2),

    reference_type VARCHAR(50),
    reference_id UUID,

    description TEXT,
    notes TEXT,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_promoter_transactions_promoter ON promoter_transaction_history(promoter_id);
CREATE INDEX idx_promoter_transactions_type ON promoter_transaction_history(transaction_type);
CREATE INDEX idx_promoter_transactions_created ON promoter_transaction_history(created_at);