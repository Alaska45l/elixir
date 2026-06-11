import { dev } from "$app/environment";
import { env as privateEnv } from "$env/dynamic/private";
import { PUBLIC_API_URL } from "$env/static/public";
import type { RequestHandler } from "./$types";

const hopByHopHeaders = new Set([
  "connection",
  "content-encoding",
  "content-length",
  "host",
  "keep-alive",
  "transfer-encoding",
  "upgrade",
]);

const responseHeaderAllowlist = new Set([
  "cache-control",
  "content-disposition",
  "content-type",
  "etag",
  "last-modified",
  "location",
  "set-cookie",
  "vary",
  "x-request-id",
]);

const proxy: RequestHandler = async ({ request, url, fetch, getClientAddress }) => {
  const configuredBackend = (privateEnv.PRIVATE_API_URL || PUBLIC_API_URL || "").replace(/\/$/, "");
  if (!configuredBackend && !dev) {
    return Response.json(
      { error: "Servicio no configurado" },
      { status: 500 },
    );
  }
  const backend = configuredBackend || "http://localhost:8080";
  const path = url.pathname.replace(/^\/api\/?/, "");
  const target = `${backend}/api/${path}${url.search}`;
  const headers = new Headers(request.headers);
  for (const header of hopByHopHeaders) headers.delete(header);
  headers.set("X-Forwarded-Proto", url.protocol.replace(":", ""));
  headers.set("X-Forwarded-Host", url.host);
  const clientAddress = getClientAddress();
  const forwardedFor = request.headers.get("x-forwarded-for");
  headers.set("X-Forwarded-For", forwardedFor ? `${forwardedFor}, ${clientAddress}` : clientAddress);
  headers.set("X-Real-IP", clientAddress);

  const method = request.method.toUpperCase();
  const body =
    method === "GET" || method === "HEAD"
      ? undefined
      : await request.arrayBuffer();
  let upstream: Response;
  try {
    upstream = await fetch(target, { method, headers, body });
  } catch {
    return Response.json(
      { error: "Servicio temporalmente no disponible" },
      { status: 502 },
    );
  }
  const responseHeaders = new Headers(upstream.headers);
  for (const header of Array.from(responseHeaders.keys())) {
    if (
      hopByHopHeaders.has(header.toLowerCase()) ||
      !responseHeaderAllowlist.has(header.toLowerCase())
    ) {
      responseHeaders.delete(header);
    }
  }

  if (upstream.status >= 500) {
    const requestID = upstream.headers.get("x-request-id");
    const payload = requestID
      ? { error: "Servicio temporalmente no disponible", request_id: requestID }
      : { error: "Servicio temporalmente no disponible" };
    return Response.json(payload, {
      status: upstream.status,
      headers: responseHeaders,
    });
  }

  return new Response(upstream.body, {
    status: upstream.status,
    headers: responseHeaders,
  });
};

export const GET = proxy;
export const POST = proxy;
export const PUT = proxy;
export const DELETE = proxy;
export const OPTIONS = proxy;
