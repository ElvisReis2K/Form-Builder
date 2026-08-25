import { apiRequest } from '../../../lib/api';
import type { BuilderForm, FormListResponse, FormRequest } from '../types';

export function listForms() {
  return apiRequest<FormListResponse>('/api/forms');
}

export function createForm(input: FormRequest) {
  return apiRequest<BuilderForm>('/api/forms', {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export function updateForm(formId: string, input: FormRequest) {
  return apiRequest<BuilderForm>(`/api/forms/${formId}`, {
    method: 'PUT',
    body: JSON.stringify(input),
  });
}

export function deleteForm(formId: string) {
  return apiRequest<void>(`/api/forms/${formId}`, {
    method: 'DELETE',
  });
}

export function publishForm(formId: string) {
  return apiRequest<BuilderForm>(`/api/forms/${formId}/publish`, {
    method: 'POST',
  });
}

export function unpublishForm(formId: string) {
  return apiRequest<BuilderForm>(`/api/forms/${formId}/unpublish`, {
    method: 'POST',
  });
}

export function getPublishedForm(slug: string) {
  return apiRequest<BuilderForm>(`/api/public/forms/${slug}`);
}
