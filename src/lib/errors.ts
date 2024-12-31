export const API_ERROR_CODES = {
  // AUTO-ADD
  INVALID_CREDENTIALS: 'invalid_credentials',
  BAD_REQUEST: 'bad_request',
  BAD_JWT: 'bad_jwt',
  UNAUTHORIZED: 'unauthorized',
  INTERNAL_SERVER_ERROR: 'internal_server_error',
} as const;

export type APIErrorCode =
  (typeof API_ERROR_CODES)[keyof typeof API_ERROR_CODES];

export class APIError extends Error {
  errorCode: APIErrorCode;
  /**
   * @param errorCode
   */
  constructor(message: string, errorCode: APIErrorCode) {
    super(message);
    this.name = 'APIError';
    this.errorCode = errorCode;
  }
}
