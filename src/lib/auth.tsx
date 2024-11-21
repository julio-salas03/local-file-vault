import {
  Component,
  createContext,
  createEffect,
  JSXElement,
  useContext,
} from 'solid-js';
import { z } from 'zod';
import { createStore, SetStoreFunction } from 'solid-js/store';
import { ServerResponseSchema } from './utils';
import { APIError } from './errors';

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
    try {
      const response = await fetch('/api/user/auth');
      const _data = await response.json();
      const schema = ServerResponseSchema(AUTH_USER_SCHEMA);
      const parse = schema.parse(_data);

      if (parse.type === 'error') {
        throw new APIError(parse.message, parse.errorCode);
      }

      setAuth('authUser', parse.data);
    } catch (error) {
      // Temporally log the error
      console.error(error);
    }
  });

  return (
    <AuthContext.Provider value={{ auth, setAuth }}>
      {props.children}
    </AuthContext.Provider>
  );
};
