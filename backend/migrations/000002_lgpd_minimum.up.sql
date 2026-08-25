ALTER TABLE forms
  ADD COLUMN controller_email text,
  ADD COLUMN privacy_purpose text,
  ADD COLUMN retention_policy text;

ALTER TABLE form_responses
  ADD COLUMN privacy_acknowledged_at timestamptz;

UPDATE form_responses
SET privacy_acknowledged_at = submitted_at
WHERE privacy_acknowledged_at IS NULL;

ALTER TABLE form_responses
  ALTER COLUMN privacy_acknowledged_at SET NOT NULL;
