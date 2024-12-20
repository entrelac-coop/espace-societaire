import { writable, derived } from "svelte/store";
import { router } from "tinro";
import jwtDecode from "jwt-decode";
import * as Sentry from "@sentry/browser";

interface JwtPayload {
  user_id: string;
  admin: string;
}

let currentToken: string | null;

export const token = writable<string | null>(localStorage.getItem("token"));
const decodedToken = derived(token, ($token) =>
  $token ? jwtDecode<JwtPayload>($token) : null
);
export const isAdmin = derived(
  decodedToken,
  ($decodedToken) => $decodedToken?.admin
);

token.subscribe((value) => {
  if (value === null) {
    localStorage.removeItem("token");
  } else {
    localStorage.setItem("token", value);
  }

  currentToken = value;
});

if (import.meta.env.PROD) {
  decodedToken.subscribe((value) => {
    if (value === null) {
      Sentry.setUser(null);
    } else {
      Sentry.setUser({ id: value.user_id });
    }
  });
}

export function getToken(): string | null {
  return currentToken;
}

export function setToken(newToken: string) {
  token.set(newToken);
}

export function signOut() {
  token.set(null);
  router.goto("/");
}

export default {
  getToken,
  setToken,
  token,
  signOut,
};
