CREATE TABLE users (
    username TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    coins INTEGER NOT NULL DEFAULT 1000
);

CREATE TABLE coin_transfers (
    id SERIAL PRIMARY KEY,
    from_username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    to_username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    amount INTEGER NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE products (
    item TEXT PRIMARY KEY,
    price INTEGER NOT NULL
);

INSERT INTO products (item, price) VALUES
  ('t-shirt', 80),
  ('cup', 20),
  ('book', 50),
  ('pen', 10),
  ('powerbank', 200),
  ('hoody', 300),
  ('umbrella', 200),
  ('socks', 10),
  ('wallet', 50),
  ('pink-hoody', 500);

CREATE TABLE purchases (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL REFERENCES users(username) ON DELETE CASCADE,
    item TEXT NOT NULL REFERENCES products(item) ON DELETE CASCADE,
    price INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
