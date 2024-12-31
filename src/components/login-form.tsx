import { Component } from 'solid-js';
import { Button } from './ui/button';
import { Label } from './ui/label';
import { TextField, TextFieldInput } from './ui/text-field';
import { toast } from 'solid-sonner';
import { AUTH_USER_SCHEMA, useAuthContext } from '../lib/auth';
import { ServerResponseSchema } from '@/lib/utils';
import { API_ERROR_CODES, APIError } from '@/lib/errors';

const LoginForm: Component = () => {
  const { setAuth } = useAuthContext();
  return (
    <form
      onSubmit={async function (e) {
        e.preventDefault();
        try {
          const response = await fetch('/api/user/login', {
            method: 'POST',
            body: new FormData(e.currentTarget),
          });
          const data = await response.json();
          const schema = ServerResponseSchema(AUTH_USER_SCHEMA);
          const parse = schema.parse(data);

          if (parse.type === 'error') {
            throw new APIError(parse.message, parse.errorCode);
          }
          setAuth('authUser', parse.data);
          toast('Logged In!');
        } catch (error) {
          if (error instanceof APIError) {
            if (error.errorCode === API_ERROR_CODES.INVALID_CREDENTIALS) {
              toast(
                'Your username and password did not match. Verify your credentials and try again'
              );
              return;
            }
          }
          toast("An unknown error has occurred and we couldn't log you in");
        }
      }}
      class="space-y-3"
    >
      <Label class="block" for="username">
        Username
      </Label>
      <TextField>
        <TextFieldInput required id="username" name="username" type="text" />
      </TextField>
      <Label for="password" class="block">
        Password
      </Label>
      <TextField>
        <TextFieldInput
          required
          id="password"
          name="password"
          type="password"
        />
      </TextField>
      <Button type="submit">Login</Button>
    </form>
  );
};

export default LoginForm;
