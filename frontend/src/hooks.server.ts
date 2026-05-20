import type { Handle, HandleServerError } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
  const response = await resolve(event);
  response.headers.set('X-Frame-Options', 'DENY');
  response.headers.set('X-Content-Type-Options', 'nosniff');
  return response;
};

export const handleError: HandleServerError = ({ error, event }) => {
  console.error('[SSR Error]', event.url.pathname, error);
  return { message: 'Error interno del servidor.' };
};
