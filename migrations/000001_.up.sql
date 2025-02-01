CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    refresh_t TEXT,
    role VARCHAR(50) CHECK (role IN ('admin', 'user')) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE weather (
    id UUID PRIMARY key not NULL,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    lat DECIMAL(9,6) NOT NULL,
    lon DECIMAL(9,6) NOT NULL,
    temp_c DECIMAL(5,2) NOT NULL,
    temp_color VARCHAR(50) NOT NULL,
    wind_kph DECIMAL(5,2) NOT NULL,
    wind_color VARCHAR(50) NOT NULL,
    cloud VARCHAR(100) NOT NULL,
    cloud_color VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
