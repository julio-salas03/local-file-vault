import { Component } from 'solid-js';
import { Button } from './ui/button';
import { Label } from './ui/label';
import { TextField, TextFieldInput } from './ui/text-field';
import { toast } from 'solid-sonner';
import { AUTH_USER_SCHEMA, useAuthContext } from '../lib/auth';

const LoginForm: Component = () => {
  const { setAuth } = useAuthContext();
  return (
    <form
      onSubmit={async function (e) {
        e.preventDefault();
        const response = await fetch('/api/user/login', {
          method: 'POST',
          body: new FormData(e.currentTarget),
        });
        const data = await response.json();
        const parse = AUTH_USER_SCHEMA.safeParse(data);

        if (!parse.success) {
          return toast("Couldn't log you in");
        }
        setAuth('authUser', parse.data);
        toast('Logged In');
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
