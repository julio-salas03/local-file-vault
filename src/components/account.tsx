import { Show } from 'solid-js';
import { useAuthContext } from '../lib/auth';
import LoginForm from './login-form';
import { Separator } from './ui/separator';

export default function Account() {
  const { auth } = useAuthContext();
  return (
    <>
      <div class="space-y-2">
        <h2 class="text-2xl font-semibold">
          <Show when={auth.authUser} fallback={'Guest Mode'}>
            {user => `Logged in as "${user().username}"`}
          </Show>
        </h2>
        <Show when={!auth.authUser}>
          <p>You can download shared files, but you can't upload your own.</p>
          <p>
            Click <strong>here (WIP)</strong> to create an account or use the
            form below to login.
          </p>
          <Separator class="my-5" />
          <LoginForm />
        </Show>
      </div>
    </>
  );
}
