export type FormStatus = 'draft' | 'published';

export type FieldType = 'text' | 'textarea' | 'email' | 'number' | 'select' | 'checkbox';

export type FormField = {
  id: string;
  position: number;
  type: FieldType;
  label: string;
  required: boolean;
  placeholder: string | null;
  options: string[];
  config: Record<string, unknown>;
};

export type BuilderForm = {
  id: string;
  title: string;
  description: string | null;
  status: FormStatus;
  publicSlug: string | null;
  publicUrl: string | null;
  publishedAt: string | null;
  fields: FormField[];
  createdAt: string;
  updatedAt: string;
};

export type FormFieldInput = {
  type: FieldType;
  label: string;
  required: boolean;
  placeholder?: string | null;
  options?: string[];
  config?: Record<string, unknown>;
};

export type FormRequest = {
  title: string;
  description?: string | null;
  fields: FormFieldInput[];
};

export type FormListResponse = {
  forms: BuilderForm[];
};

export type FormSubmission = {
  id: string;
  formId: string;
  answers: Record<string, unknown>;
  submittedAt: string;
};

export type SubmitResponseRequest = {
  answers: Record<string, unknown>;
};

export type FormResponseSummary = {
  id: string;
  title: string;
  fields: Pick<FormField, 'id' | 'position' | 'type' | 'label' | 'required'>[];
};

export type FormSubmissionListResponse = {
  form: FormResponseSummary;
  responses: FormSubmission[];
};
