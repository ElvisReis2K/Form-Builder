ALTER TABLE form_fields
  DROP CONSTRAINT IF EXISTS form_fields_type_check;

ALTER TABLE form_fields
  ADD CONSTRAINT form_fields_type_check
  CHECK (type IN ('text', 'textarea', 'email', 'number', 'phone', 'select', 'checkbox'));
