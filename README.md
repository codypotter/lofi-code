# lofi-code

`lofi-code` is a personal blogging project built with [Angular CLI](https://github.com/angular/angular-cli).

## Getting Started

Clone this repository. Run `npm install` with a modern node version, 20 or greater. Create an `environments.dev.ts` from `environments.ts` with valid secrets. Then `npm start`.

## Server-Side Rendering (SSR)

This project uses Angular Universal for server-side rendering. This allows the application to render the initial page on the server, which can improve performance and make the application more crawlable for SEO purposes. To run the SSR server, use the command `npm run dev:ssr`.

## Running unit tests

Run `npm run test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## AWS Deployment

I'm avoiding committing the `environment.prod.ts` config, even though it's used on the client and it does not contain true secrets. Keeping it out of git reinforces good security hygiene and reduces the risk of it being scraped, reused, or misunderstood as safe to treat like a backend secret.

Trigger a deployment by creating a GitHub release from a tag.

### Local Deployment

This project is configured for deployment on AWS. To deploy the application, ensure you have a valid `environments.prod.ts` file. Then authenticate your command line with aws by adding your aws credentials to `~/aws/.credentials` and running `npm run deploy`.
