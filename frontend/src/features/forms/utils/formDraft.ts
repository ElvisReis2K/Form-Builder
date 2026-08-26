import type { BuilderForm, FieldType, FormField, FormRequest } from '../types';

export const fieldTypes: FieldType[] = ['text', 'textarea', 'email', 'number', 'phone', 'select', 'checkbox'];

export const fieldTypeLabels: Record<FieldType, string> = {
  text: 'Texto curto',
  textarea: 'Texto longo',
  email: 'E-mail',
  number: 'Número',
  phone: 'Telefone',
  select: 'Seleção',
  checkbox: 'Caixa de seleção',
};

export type FieldDraft = {
  id: string | null;
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
  controllerEmail: string;
  privacyPurpose: string;
  retentionPolicy: string;
  fields: FieldDraft[];
};

let draftSequence = 0;

export function createBlankDraft(): FormDraft {
  return {
    id: null,
    title: 'Formulário sem título',
    description: '',
    controllerEmail: '',
    privacyPurpose: '',
    retentionPolicy: '',
    fields: [createFieldDraft()],
  };
}

export function createFieldDraft(field?: Partial<FormField>): FieldDraft {
  draftSequence += 1;

  return {
    id: field?.id ?? null,
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
    controllerEmail: form.controllerEmail ?? '',
    privacyPurpose: form.privacyPurpose ?? '',
    retentionPolicy: form.retentionPolicy ?? '',
    fields: form.fields.map((field) => createFieldDraft(field)),
  };
}

export function draftToRequest(draft: FormDraft): FormRequest {
  return {
    title: draft.title,
    description: optionalText(draft.description),
    controllerEmail: optionalText(draft.controllerEmail),
    privacyPurpose: optionalText(draft.privacyPurpose),
    retentionPolicy: optionalText(draft.retentionPolicy),
    fields: draft.fields.map((field) => ({
      id: field.id ?? undefined,
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
