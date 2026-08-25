import type { BuilderForm, FieldType, FormField, FormRequest } from '../types';

export const fieldTypes: FieldType[] = ['text', 'textarea', 'email', 'number', 'select', 'checkbox'];

export const fieldTypeLabels: Record<FieldType, string> = {
  text: 'Short text',
  textarea: 'Long text',
  email: 'Email',
  number: 'Number',
  select: 'Select',
  checkbox: 'Checkbox',
};

export type FieldDraft = {
  clientId: string;
  type: FieldType;
  label: string;
  required: boolean;
  placeholder: string;
  optionsText: string;
};

export type FormDraft = {
  id: string | null;
  title: string;
  description: string;
  fields: FieldDraft[];
};

let draftSequence = 0;

export function createBlankDraft(): FormDraft {
  return {
    id: null,
    title: 'Untitled form',
    description: '',
    fields: [createFieldDraft()],
  };
}

export function createFieldDraft(field?: Partial<FormField>): FieldDraft {
  draftSequence += 1;

  return {
    clientId: field?.id ?? `new-field-${Date.now()}-${draftSequence}`,
    type: field?.type ?? 'text',
    label: field?.label ?? '',
    required: field?.required ?? false,
    placeholder: field?.placeholder ?? '',
    optionsText: field?.options?.join('\n') ?? '',
  };
}

export function formToDraft(form: BuilderForm): FormDraft {
  return {
    id: form.id,
    title: form.title,
    description: form.description ?? '',
    fields: form.fields.map((field) => createFieldDraft(field)),
  };
}

export function draftToRequest(draft: FormDraft): FormRequest {
  return {
    title: draft.title,
    description: optionalText(draft.description),
    fields: draft.fields.map((field) => ({
      type: field.type,
      label: field.label,
      required: field.required,
      placeholder: optionalText(field.placeholder),
      options: normalizeOptions(field.optionsText),
      config: {},
    })),
  };
}

function optionalText(value: string) {
  const trimmed = value.trim();
  return trimmed === '' ? null : trimmed;
}

function normalizeOptions(value: string) {
  const seen = new Set<string>();

  return value
    .split('\n')
    .map((option) => option.trim())
    .filter((option) => {
      if (option === '' || seen.has(option)) {
        return false;
      }

      seen.add(option);
      return true;
    });
}
