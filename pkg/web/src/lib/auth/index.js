// @flow

const key = "auth/session";

export type Session = {
  userUUID: string,
  sessionToken: string,
};

export function storeSession(userUUID: string, sessionToken: string) {
  window.localStorage.setItem(
    key,
    JSON.stringify({
      userUUID,
      sessionToken,
    })
  );
}

export function getSession(): ?Session {
  const item = window.localStorage.getItem(key);

  return item ? JSON.parse(item) : null;
}

export function clearSession() {
  window.localStorage.removeItem(key);
}
