import auth from "./auth";

export let baseURL = "http://localhost:8080/";

if (import.meta.env.PROD) {
  baseURL = "https://societaire.entrelac.coop/api/";
}

function camelCasedKey(key: string): string {
  return key.replace(/([A-Z])/g, "_$1").toLowerCase();
}

function camelCasedObject(object: any): any {
  const newObject: any = {};

  for (const key in object) {
    newObject[camelCasedKey(key)] = object[key];
  }

  return newObject;
}

async function call(
  method: string,
  path: string,
  headers: any,
  body?: string | FormData
) {
  const currentToken = auth.getToken();

  if (currentToken) {
    headers = {
      ...headers,
      Authorization: "Bearer " + currentToken,
    };
  }

  const response = await fetch(baseURL + path, {
    method,
    cache: "no-cache",
    headers,
    body,
  });

  const json = await response.json();

  if (!response.ok) {
    if (json.code === "token-expired" || json.code === "token-invalid") {
      auth.signOut();
    }

    throw json;
  }

  return json;
}

async function post(path: string, data: any) {
  const headers = {
    "Content-Type": "application/json",
  };

  const body = JSON.stringify(camelCasedObject(data));

  return await call("POST", path, headers, body);
}

async function upload(path: string, body: FormData) {
  const headers = {};

  return await call("POST", path, headers, body);
}

async function get(path: string) {
  const headers = {
    "Content-Type": "application/json",
  };

  return await call("GET", path, headers);
}

interface CreateUserRequest {
  email: string;
  password: string;
  phoneNumber: string;
  firstName: string;
  lastName: string;
  address: string;
  postalCode: string;
  city: string;
  country: string;
  category: string;
  reason: string | null;
}

export async function createUser(user: CreateUserRequest) {
  return await post("users", user);
}

interface ConfirmUserRequest {
  email: string;
  token: string;
}

export async function confirmUser(data: ConfirmUserRequest) {
  return await post("users/confirm", data);
}

interface StartResetUserRequest {
  email: string;
}

export async function startResetUser(data: StartResetUserRequest) {
  return await post("users/reset/start", data);
}

interface ResetUserRequest {
  email: string;
  password: string;
  token: string;
}

export async function resetUser(data: ResetUserRequest) {
  return await post("users/reset", data);
}

interface StartConfirmUserRequest {
  email: string;
}

export async function startConfirmUser(data: StartConfirmUserRequest) {
  return await post("users/confirm/start", data);
}

interface CreateTokenRequest {
  email: string;
  password: string;
}

export async function createToken(credentials: CreateTokenRequest) {
  return await post("tokens", credentials);
}

interface CreateCheckoutSessionRequest {
  quantity: number;
  gift?: boolean;
}

export async function createCheckoutSession(
  data: CreateCheckoutSessionRequest
) {
  return await post("users/me/checkout/sessions", data);
}

interface UseGiftCodeRequest {
  giftCode: string;
}

export async function useGiftCode(data: UseGiftCodeRequest) {
  return await post("users/me/use-gift-code", data);
}

export async function getCurrentUser() {
  return await get("users/me");
}

export async function getUser(userID: string) {
  return await get(`admin/users/${userID}`);
}

export async function getUsers() {
  return await get("admin/users");
}

export async function uploadDocuments(body: FormData) {
  return await upload("users/me/documents", body);
}
