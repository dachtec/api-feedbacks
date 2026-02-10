CREATE TABLE IF NOT EXISTS feedbacks (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       VARCHAR(255) NOT NULL,
    feedback_type VARCHAR(50)  NOT NULL CHECK (feedback_type IN ('bug','suggestion','praise','question')),
    rating        INTEGER      NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment       TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_user_id    ON feedbacks(user_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_type       ON feedbacks(feedback_type);
CREATE INDEX IF NOT EXISTS idx_feedbacks_rating     ON feedbacks(rating);
CREATE INDEX IF NOT EXISTS idx_feedbacks_created_at ON feedbacks(created_at);
