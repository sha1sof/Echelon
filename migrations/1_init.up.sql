CREATE TABLE IF NOT EXISTS video
(
    id INTEGER PRIMARY KEY,
    video_id TEXT NOT NULL UNIQUE,
    hash BLOB
);
CREATE INDEX IF NOT EXISTS idx_video_id ON video (video_id);