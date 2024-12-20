import "./app.css";
import App from "./App.svelte";

import * as Sentry from "@sentry/svelte";
import { BrowserTracing } from "@sentry/tracing";

if (import.meta.env.PROD) {
  Sentry.init({
    dsn: "https://bef41888eb434dc2b6b3af14687a70fa@o1422105.ingest.sentry.io/6768664",
    integrations: [new BrowserTracing()],
    tracesSampleRate: 1.0,
  });
}

const app = new App({
  target: document.getElementById("app")!,
});

export default app;
