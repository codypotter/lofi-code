const { default: serverlessExpress } = await import('@codegenie/serverless-express');
import app from './dist/lofi-code/server/server.mjs';



export const handler = serverlessExpress({
  app,
  binarySettings: {
    contentTypes: [
      "application/javascript",
      "application/json",
      "application/octet-stream",
      "application/xml",
      "image/jpeg",
      "image/png",
      "image/gif",
      "text/comma-separated-values",
      "text/css",
      "text/html",
      "text/javascript",
      "text/plain",
      "text/text",
      "text/xml",
      "image/x-icon",
      "image/svg+xml",
      "application/x-font-ttf",
    ]
  }
});
