import {
  Component,
  createContext,
  createEffect,
  JSXElement,
  useContext,
} from 'solid-js';
import { z } from 'zod';
import { createStore, SetStoreFunction } from 'solid-js/store';

export const AUTH_USER_SCHEMA = z.object({
  username: z.string(),
  uploadFolder: z.string(),
});

type AuthUser = z.infer<typeof AUTH_USER_SCHEMA> | null;

type AuthContextProps = {
  authUser: AuthUser;
};

const AuthContext = createContext({
  auth: { authUser: null } as AuthContextProps,
  setAuth: function (...args: any) {} as SetStoreFunction<AuthContextProps>,
});

export function useAuthContext() {
  return useContext(AuthContext);
}

export const AuthProvider: Component<{ children: JSXElement }> = props => {
  const [auth, setAuth] = createStore<AuthContextProps>({
    authUser: null,
  });

  createEffect(async () => {
    const response = await fetch('/api/user/auth');
    const data = await response.json();
    const parse = AUTH_USER_SCHEMA.safeParse(data);

    if (!parse.success) {
      console.log(parse.error);
      return;
    }
    setAuth('authUser', parse.data);
  });

  return (
    <AuthContext.Provider value={{ auth, setAuth }}>
      {props.children}
    </AuthContext.Provider>
  );
};
