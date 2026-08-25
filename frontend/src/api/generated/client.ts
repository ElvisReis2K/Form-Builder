/**
 * This file was auto-generated from backend/openapi/openapi.json.
 * Do not edit manually. Run npm run generate:api from frontend instead.
 */
import type { components, paths } from './schema';

const API_BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

type APIErrorPayload = Partial<components['schemas']['ErrorResponse']>;
type HTTPMethod = 'get' | 'post' | 'put' | 'delete' | 'patch';
type PathItem<TPath extends keyof paths> = paths[TPath];
type AvailableMethod<TPath extends keyof paths> = {
  [TMethod in HTTPMethod]: TMethod extends keyof PathItem<TPath>
    ? Exclude<PathItem<TPath>[TMethod], undefined> extends never
      ? never
      : TMethod
    : never;
}[HTTPMethod];
type Operation<TPath extends keyof paths, TMethod extends AvailableMethod<TPath>> = TMethod extends keyof PathItem<TPath>
  ? Exclude<PathItem<TPath>[TMethod], undefined>
  : never;
type JSONRequestBody<TOperation> = TOperation extends {
  requestBody: {
    content: {
      'application/json': infer TBody;
    };
  };
}
  ? TBody
  : never;
type PathParameters<TOperation> = TOperation extends {
  parameters: {
    path: infer TPath;
  };
}
  ? TPath
  : never;
type JSONResponse<TResponse> = TResponse extends {
  content: {
    'application/json': infer TBody;
  };
}
  ? TBody
  : void;
type StatusResponse<TResponses, TStatus extends number> = TStatus extends keyof TResponses
  ? JSONResponse<TResponses[TStatus]>
  : never;
type SuccessResponse<TOperation> = TOperation extends {
  responses: infer TResponses;
}
  ? StatusResponse<TResponses, 200> | StatusResponse<TResponses, 201> | StatusResponse<TResponses, 204>
  : never;
type PathOption<TOperation> = [PathParameters<TOperation>] extends [never]
  ? { path?: never }
  : { path: PathParameters<TOperation> };
type BodyOption<TOperation> = [JSONRequestBody<TOperation>] extends [never]
  ? { body?: never }
  : { body: JSONRequestBody<TOperation> };
type APIRequestOptions<TOperation> = PathOption<TOperation> &
  BodyOption<TOperation> & {
    request?: Omit<RequestInit, 'body' | 'method'>;
  };
type RequiresOptions<TOperation> = [PathParameters<TOperation>] extends [never]
  ? [JSONRequestBody<TOperation>] extends [never]
    ? false
    : true
  : true;
type APIRequestArgs<TOperation> = RequiresOptions<TOperation> extends true
  ? [options: APIRequestOptions<TOperation>]
  : [options?: APIRequestOptions<TOperation>];
type RuntimePathParameters = Record<string, string | number | boolean | null | undefined>;
type RuntimeAPIRequestOptions = {
  path?: RuntimePathParameters;
  body?: unknown;
  request?: Omit<RequestInit, 'body' | 'method'>;
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

  return 'Erro inesperado na requisicao';
}

export function apiURL<TPath extends keyof paths>(path: TPath, pathParams?: RuntimePathParameters) {
  return absoluteAPIURL(resolvePath(String(path), pathParams));
}

export function deleteApiFormsFormId(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}", "delete">>
) {
  return apiRequest("/api/forms/{formId}", "delete", ...args);
}

export function deleteApiFormsFormIdResponsesResponseId(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}/responses/{responseId}", "delete">>
) {
  return apiRequest("/api/forms/{formId}/responses/{responseId}", "delete", ...args);
}

export function getApiAuthGoogle(
  ...args: APIRequestArgs<Operation<"/api/auth/google", "get">>
) {
  return apiRequest("/api/auth/google", "get", ...args);
}

export function getApiAuthGoogleCallback(
  ...args: APIRequestArgs<Operation<"/api/auth/google/callback", "get">>
) {
  return apiRequest("/api/auth/google/callback", "get", ...args);
}

export function getApiAuthMe(
  ...args: APIRequestArgs<Operation<"/api/auth/me", "get">>
) {
  return apiRequest("/api/auth/me", "get", ...args);
}

export function getApiForms(
  ...args: APIRequestArgs<Operation<"/api/forms", "get">>
) {
  return apiRequest("/api/forms", "get", ...args);
}

export function getApiFormsFormId(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}", "get">>
) {
  return apiRequest("/api/forms/{formId}", "get", ...args);
}

export function getApiFormsFormIdResponses(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}/responses", "get">>
) {
  return apiRequest("/api/forms/{formId}/responses", "get", ...args);
}

export function getApiFormsFormIdResponsesExport(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}/responses/export", "get">>
) {
  return apiRequest("/api/forms/{formId}/responses/export", "get", ...args);
}

export function getApiPublicFormsSlug(
  ...args: APIRequestArgs<Operation<"/api/public/forms/{slug}", "get">>
) {
  return apiRequest("/api/public/forms/{slug}", "get", ...args);
}

export function getHealthz(
  ...args: APIRequestArgs<Operation<"/healthz", "get">>
) {
  return apiRequest("/healthz", "get", ...args);
}

export function getReadyz(
  ...args: APIRequestArgs<Operation<"/readyz", "get">>
) {
  return apiRequest("/readyz", "get", ...args);
}

export function postApiAuthLogin(
  ...args: APIRequestArgs<Operation<"/api/auth/login", "post">>
) {
  return apiRequest("/api/auth/login", "post", ...args);
}

export function postApiAuthLogout(
  ...args: APIRequestArgs<Operation<"/api/auth/logout", "post">>
) {
  return apiRequest("/api/auth/logout", "post", ...args);
}

export function postApiAuthRegister(
  ...args: APIRequestArgs<Operation<"/api/auth/register", "post">>
) {
  return apiRequest("/api/auth/register", "post", ...args);
}

export function postApiForms(
  ...args: APIRequestArgs<Operation<"/api/forms", "post">>
) {
  return apiRequest("/api/forms", "post", ...args);
}

export function postApiFormsFormIdPublish(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}/publish", "post">>
) {
  return apiRequest("/api/forms/{formId}/publish", "post", ...args);
}

export function postApiFormsFormIdUnpublish(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}/unpublish", "post">>
) {
  return apiRequest("/api/forms/{formId}/unpublish", "post", ...args);
}

export function postApiPublicFormsSlugResponses(
  ...args: APIRequestArgs<Operation<"/api/public/forms/{slug}/responses", "post">>
) {
  return apiRequest("/api/public/forms/{slug}/responses", "post", ...args);
}

export function putApiFormsFormId(
  ...args: APIRequestArgs<Operation<"/api/forms/{formId}", "put">>
) {
  return apiRequest("/api/forms/{formId}", "put", ...args);
}

async function apiRequest<TPath extends keyof paths, TMethod extends AvailableMethod<TPath>>(
  path: TPath,
  method: TMethod,
  ...args: APIRequestArgs<Operation<TPath, TMethod>>
): Promise<SuccessResponse<Operation<TPath, TMethod>>> {
  const options = (args[0] ?? {}) as RuntimeAPIRequestOptions;
  const request = options.request ?? {};
  const headers = new Headers(request.headers);
  if (options.body !== undefined && !headers.has('Content-Type')) {
    headers.set('Content-Type', 'application/json');
  }

  const response = await fetch(absoluteAPIURL(resolvePath(String(path), options.path)), {
    ...request,
    method: method.toUpperCase(),
    credentials: 'include',
    headers,
    body: options.body === undefined ? undefined : JSON.stringify(options.body),
  });

  if (response.status === 204) {
    return undefined as SuccessResponse<Operation<TPath, TMethod>>;
  }

  const payload = await readPayload(response);
  if (!response.ok) {
    const apiError = payload as APIErrorPayload | undefined;
    throw new APIError(
      apiError?.error?.message ?? 'Falha na requisicao',
      response.status,
      apiError?.error?.code ?? 'request_failed',
    );
  }

  return payload as SuccessResponse<Operation<TPath, TMethod>>;
}

function readPayload(response: Response) {
  const contentType = response.headers.get('content-type') ?? '';
  if (!contentType.includes('application/json')) {
    return response.text();
  }

  return response.json();
}

function absoluteAPIURL(route: string) {
  return `${API_BASE_URL}${route}`;
}

function resolvePath(route: string, pathParams: RuntimePathParameters = {}) {
  return route.replace(/\{([^}]+)\}/g, (_, key: string) => {
    const value = pathParams[key];
    if (value === undefined || value === null) {
      throw new Error(`Parametro de rota ausente: "${key}"`);
    }

    return encodeURIComponent(String(value));
  });
}
