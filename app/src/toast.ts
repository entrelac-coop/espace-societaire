import { writable } from "svelte/store";

interface Toast {
  type: string;
  message: string;
}

export const currentToast = writable<Toast | null>(null);

function set(toast: Toast) {
  currentToast.set(toast);

  setTimeout(() => {
    currentToast.update((value) => {
      if (value === toast) {
        return null;
      }

      return value;
    });
  }, 1000 * 5);
}

function error(message: string) {
  set({ type: "error", message });
}

function success(message: string) {
  set({ type: "success", message });
}

export default {
  set,
  error,
  success,
};
