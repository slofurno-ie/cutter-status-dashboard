CREATE TABLE statuses (
    status_id SERIAL PRIMARY KEY,
    service text,
    status text
);

INSERT INTO statuses (service) VALUES
('platform'),
('fulfillment'),
('crm'),
('study')
;