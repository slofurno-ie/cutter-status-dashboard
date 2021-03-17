CREATE TABLE statuses (
    status_id SERIAL PRIMARY KEY,
    service text,
    status text DEFAULT '200'
);

INSERT INTO statuses (service) VALUES
('platform'),
('fulfillment'),
('crm'),
('study')
;