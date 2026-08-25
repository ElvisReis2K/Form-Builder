import { apiRequest } from '../../../lib/api';
import type { FormRequest, SubmitResponseRequest } from '../types';

export function listForms() {
  return apiRequest('/api/forms', 'get');
}

export function createForm(input: FormRequest) {
  return apiRequest('/api/forms', 'post', {
    body: input,
  });
}

export function updateForm(formId: string, input: FormRequest) {
  return apiRequest('/api/forms/{formId}', 'put', {
    path: { formId },
    body: input,
  });
}

export function deleteForm(formId: string) {
  return apiRequest('/api/forms/{formId}', 'delete', {
    path: { formId },
  });
}

export function publishForm(formId: string) {
  return apiRequest('/api/forms/{formId}/publish', 'post', {
    path: { formId },
  });
}

export function unpublishForm(formId: string) {
  return apiRequest('/api/forms/{formId}/unpublish', 'post', {
    path: { formId },
  });
}

export function getPublishedForm(slug: string) {
  return apiRequest('/api/public/forms/{slug}', 'get', {
    path: { slug },
  });
}

export function submitFormResponse(slug: string, input: SubmitResponseRequest) {
  return apiRequest('/api/public/forms/{slug}/responses', 'post', {
    path: { slug },
    body: input,
  });
}

export function listFormResponses(formId: string) {
  return apiRequest('/api/forms/{formId}/responses', 'get', {
    path: { formId },
  });
}
