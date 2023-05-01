CREATE TABLE transcription_requests (
    wa_id TEXT PRIMARY KEY,
    from_number TEXT NOT NULL,
    to_number TEXT NOT NULL,
    media_url TEXT NOT NULL,
    received_time TIMESTAMP NOT NULL,
    replied_time TIMESTAMP NOT NULL
);
