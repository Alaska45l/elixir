import type { Handle, HandleServerError } from "@sveltejs/kit";
import { dev } from "$app/environment";

const csp = [
  "default-src 'self'",
  "base-uri 'self'",
  "frame-ancestors 'none'",
  "object-src 'none'",
  "img-src 'self' data: blob: https:",
  "font-src 'self' https://fonts.gstatic.com data:",
  "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
  "script-src 'self'",
  "connect-src 'self'",
  "form-action 'self'",
  "frame-src 'self' https://www.mercadopago.com https://*.mercadopago.com",
].join("; ");

export const handle: Handle = async ({ event, resolve }) => {
  const response = await resolve(event);
  response.headers.set("X-Frame-Options", "DENY");
  response.headers.set("X-Content-Type-Options", "nosniff");
  response.headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  response.headers.set("Content-Security-Policy", csp);
  response.headers.set(
    "Permissions-Policy",
    "camera=(), microphone=(), geolocation=()",
  );
  if (!dev) {
    response.headers.set(
      "Strict-Transport-Security",
      "max-age=31536000; includeSubDomains; preload",
    );
  }
  return response;
};

export const handleError: HandleServerError = ({ error, event }) => {
  console.error("[SSR Error]", event.url.pathname, error);
  return { message: "Error interno del servidor." };
};
