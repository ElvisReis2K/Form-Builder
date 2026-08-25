const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

type APIErrorPayload = {
  error?: {
    code?: string;
    message?: string;
  };
};

export class APIError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(message: string, status: number, code: string) {
    super(message);
    this.name = 'APIError';
    this.status = status;
    this.code = code;
  }
}

export function isAPIError(error: unknown): error is APIError {
  return error instanceof APIError;
}

export function getErrorMessage(error: unknown) {
  if (isAPIError(error)) {
    return error.message;
  }

  if (error instanceof Error) {
    return error.message;
  }

  return 'Unexpected request error';
}

export async function apiRequest<TResponse>(path: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers);
  if (options.body && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    credentials: 'include',
    headers,
  });

  if (response.status === 204) {
    return undefined as TResponse;
  }

  const payload = await readPayload(response);
  if (!response.ok) {
    const apiError = payload as APIErrorPayload | undefined;
    throw new APIError(
      apiError?.error?.message ?? 'Request failed',
      response.status,
      apiError?.error?.code ?? 'request_failed',
    );
  }

  return payload as TResponse;
}

function readPayload(response: Response) {
  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    return response.text();
  }

  return response.json();
}
