CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    date TEXT NOT NULL,
    amount FLOAT NOT NULL,
    type TEXT CHECK (type IN ('Credit', 'Debit')) NOT NULL
);

INSERT INTO transactions (date, amount, type) VALUES 
    ('2024-02-24', 1500, 'Credit'),
    ('2024-02-23', -500, 'Debit'),
    ('2024-02-22', 2000, 'Credit');
