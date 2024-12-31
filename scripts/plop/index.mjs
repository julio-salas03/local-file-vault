import { exec } from 'child_process';
import path from 'path';
import process from 'process';

const BASE_PATH = process.cwd();

export default function (
  /** @type {import('plop').NodePlopAPI} */
  plop
) {
  plop.setActionType('format', (answers, config) => {
    if (typeof config.file !== 'string') return;

    if (config.file.endsWith('.go')) {
      exec(`go fmt ${config.file}`);
      return;
    }
  });

  const GO_ERROR_CODES_DECLARATION_FILE = path.join(
    BASE_PATH,
    'src/server/errorcodes/index.go'
  );

  const JAVASCRIPT_ERROR_CODES_DECLARATION_FILE = path.join(
    BASE_PATH,
    'src/lib/errors.ts'
  );

  plop.setGenerator('errorcode', {
    description: 'add a new api error code',
    prompts: [{ type: 'input', name: 'name', message: 'error code name' }],
    actions: [
      {
        pattern: '// AUTO-ADD',
        type: 'append',
        path: GO_ERROR_CODES_DECLARATION_FILE,
        template: `{{ pascalCase name }} = "{{ snakeCase name }}"`,
      },

      {
        pattern: '// AUTO-ADD',
        type: 'append',
        path: JAVASCRIPT_ERROR_CODES_DECLARATION_FILE,
        template: `  {{ constantCase name }}: '{{ snakeCase name }}',`,
      },
      {
        type: 'format',
        file: GO_ERROR_CODES_DECLARATION_FILE,
      },
    ],
  });
}
