import { Component } from 'solid-js';
import { Button } from './ui/button';
import { Label } from './ui/label';
import { TextField, TextFieldInput } from './ui/text-field';
import { toast } from 'solid-sonner';

const LoginForm: Component = () => {
  return (
    <form
      onSubmit={async e => {
        e.preventDefault();
        const data = new FormData(e.currentTarget);
        console.log(data.get('username'));
        const response = await fetch('/api/login', {
          method: 'POST',
          body: new FormData(e.currentTarget),
        });
        const text = await response.text();
        toast(text);
        //window.location.reload();
      }}
      class="space-y-3"
    >
      <Label class="block" for="username">
        Username
      </Label>
      <TextField>
        <TextFieldInput id="username" name="username" type="text" />
      </TextField>
      <Label for="password" class="block">
        Password
      </Label>
      <TextField>
        <TextFieldInput id="password" name="password" type="password" />
      </TextField>
      <Button type="submit">Login</Button>
    </form>
  );
};

export default LoginForm;
