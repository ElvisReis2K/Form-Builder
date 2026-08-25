ALTER TABLE form_responses
  DROP COLUMN IF EXISTS privacy_acknowledged_at;

ALTER TABLE forms
  DROP COLUMN IF EXISTS retention_policy,
  DROP COLUMN IF EXISTS privacy_purpose,
  DROP COLUMN IF EXISTS controller_email;
