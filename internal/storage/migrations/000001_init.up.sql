CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        login VARCHAR(100) NOT NULL,
        password VARCHAR(100) NOT NULL,
        balance DECIMAL(10, 2) DEFAULT 0,
        withdrawn DECIMAL(10, 2) DEFAULT 0,
        created_at TIMESTAMP DEFAULT now () NOT NULL,
        updated_at TIMESTAMP DEFAULT now () NOT NULL,
        CONSTRAINT login_unique UNIQUE (login)
    );

CREATE TABLE
    IF NOT EXISTS orders (
        id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        user_id INTEGER NOT NULL,
        number VARCHAR(255) NOT NULL,
        status order_status NOT NULL,
        accrual DECIMAL(10, 2) DEFAULT 0,
        created_at TIMESTAMP DEFAULT now () NOT NULL,
        updated_at TIMESTAMP DEFAULT now () NOT NULL,
        CONSTRAINT number_unique unique (number),
        CONSTRAINT orders_fk_users foreign key (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS withdrawals (
        id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        user_id INTEGER NOT NULL,
        order_number VARCHAR(255) NOT NULL,
        amount DECIMAL(10, 2) NOT NULL,
        created_at TIMESTAMP DEFAULT now () NOT NULL,
        CONSTRAINT order_number_unique unique (order_number),
        CONSTRAINT withdrawals_fk_users foreign key (user_id) REFERENCES users (id)
    );