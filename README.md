# lofi-code

**lofi-code** is my personal blogging project built with a custom static site generator in Go. The blog uses progressive enhancement with templ and htmx to deliver dynamic HTML snippets via AWS Lambda and API Gateway while serving pre-built static assets from S3 through CloudFront.

---

## Getting Started

1. **Clone the Repository**  
   ```bash
   git clone https://github.com/yourusername/lofi-code.git
   cd lofi-code
   ```

2. **Local Development**  
   - To work locally, install [Go](https://golang.org/dl/) and [Docker](https://www.docker.com/).
   - Run the development loop using:
     ```bash
     make dev
     ```
     This command watches for changes in templates (via `templ`) and re-generates the static site before launching the local server. The static site is served from the `public/` directory.

3. **Static Site Generation**  
    The `generate` target runs template generation (via `templ generate`) and then uses `go run` to build and run the generator.  
    ```bash
    make generate
    ```

---

## Deployment

Your project is deployed to AWS using a combination of AWS Lambda, API Gateway, S3 (for static assets), and CloudFront for caching. The full deployment process is orchestrated via Makefile targets and a GitHub Actions workflow.

To trigger a full deployment manually:
```bash
make deploy
```

---

## AWS Deployment Details

- **AWS Lambda & API Gateway:**  
  Dynamically delivered content (e.g. htmx partials) is served via Lambda functions.
- **S3 & CloudFront:**  
  All static assets including HTML files, CSS, JS, fonts, and images are stored in S3 and cached globally using CloudFront for optimal performance.
- **IAM & CI/CD:**  
  GitHub Actions handles releases and deployments via CloudFormation, ECR, and S3, leveraging a tightly controlled IAM policy.
---

## Infrastructure Diagram

Below is the architecture diagram for the lofi-code stack:

![Infrastructure Diagram](./infra.png)

---

## Contributing

For any questions, contact me at [me@codypotter.com](mailto:me@codypotter.com).

---
