# lofi-code

`lofi-code` is a personal blogging project built with [Angular CLI](https://github.com/angular/angular-cli).

## Getting Started

Clone this repository. Run `npm install` with a modern node version, 18 or greater. Create an `environments.dev.ts` from `environments.ts` with valid secrets. Then `npm start`.

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The application will automatically reload if you change any of the source files.

## Server-Side Rendering (SSR)

This project uses Angular Universal for server-side rendering. This allows the application to render the initial page on the server, which can improve performance and make the application more crawlable for SEO purposes. To run the SSR server, use the command `npm run dev:ssr`.

## Build

Run `npm run buld` to build the project. The build artifacts will be stored in the `dist/` directory.

## Running unit tests

Run `npm run test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## AWS Deployment

This project is configured for deployment on AWS. To deploy the application, ensure you have a valid `environments.prod.ts` file. Then authenticate your command line with aws by adding your aws credentials to `~/aws/.credentials` and running `npm run deploy`.

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.