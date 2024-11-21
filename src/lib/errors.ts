export class APIError extends Error {
  errorCode: string;
  /**
   * @param errorCode NEED TO TYPE THIS WITH AN ENUM OR CONSTANTS
   */
  constructor(message: string, errorCode: string) {
    super(message);
    this.name = 'APIError';
    this.errorCode = errorCode;
  }
}
