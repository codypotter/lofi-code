import awsServerlessExpress from 'aws-serverless-express';
import awsServerlessExpressMiddleware from 'aws-serverless-express/middleware';
import app from './dist/lofi-code/server/server.mjs';

const binaryMimeTypes = [
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
];

app.use(awsServerlessExpressMiddleware.eventContext());

const server = awsServerlessExpress.createServer(app, null, binaryMimeTypes);

export const universal = (event, context) =>
  awsServerlessExpress.proxy(server, event, context);