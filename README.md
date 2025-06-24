# SMTP Form Handler

A simple Go application that receives form submissions via HTTP and sends emails using SMTP. Built to run on any platform including Cloudflare Containers.

## Features

- Receives JSON form submissions at `/send-email` endpoint
- Sends emails via SMTP using Go's standard library
- API key authentication for security
- Health check endpoint at `/health`
- Optimized Docker image for cloud deployment
- CORS support for web applications

## Environment Variables

The following environment variables are required:

- `SMTP_HOST`: SMTP server hostname (e.g., smtp.gmail.com)
- `SMTP_PORT`: SMTP server port (default: 587)
- `SMTP_USERNAME`: Email account username
- `SMTP_PASSWORD`: Email account password
- `RECIPIENT_EMAIL`: Email address where form submissions will be sent
- `API_KEY`: Secret key for API authentication
- `SERVER_PORT`: HTTP server port (default: 8080)

## Local Development

1. Copy `.env.example` to `.env` and fill in your credentials
2. Run the application:
   ```bash
   export $(cat .env | grep -v '^#' | xargs) && go run main.go
   ```

## API Usage

Send a POST request to `/send-email` with the following JSON payload:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "subject": "Contact Form",
  "message": "Your message here"
}
```

Include the API key in the header:
```
X-API-Key: your-api-key-here
```

## Docker Deployment

Build and run the Docker image:

```bash
# Build the image
docker build -t smtp-form-handler .

# Run the container
docker run -p 8080:8080 \
  -e SMTP_HOST=smtp.gmail.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=your-email@gmail.com \
  -e SMTP_PASSWORD=your-app-password \
  -e RECIPIENT_EMAIL=recipient@example.com \
  -e API_KEY=your-secret-api-key \
  smtp-form-handler
```

## Cloudflare Containers Deployment

This application is optimized for Cloudflare Containers with:
- Multi-stage build for minimal image size
- Scratch base image for security
- linux/amd64 architecture
- Port 8080 default configuration

### Prerequisites

1. Install [Wrangler CLI](https://developers.cloudflare.com/workers/wrangler/install-and-update/)
2. Authenticate with Cloudflare: `wrangler login`

### Deployment Steps

1. Create a `wrangler.toml` file in your project root:

```toml
name = "smtp-form-handler"
main = "src/index.ts"
compatibility_date = "2024-06-24"

[[containers]]
name = "smtp-backend"
image = { dockerfile = "Dockerfile" }

[env.production.vars]
SERVER_PORT = "8080"

# Secrets should be set using wrangler CLI:
# wrangler secret put SMTP_HOST
# wrangler secret put SMTP_PORT
# wrangler secret put SMTP_USERNAME
# wrangler secret put SMTP_PASSWORD
# wrangler secret put RECIPIENT_EMAIL
# wrangler secret put API_KEY
```

2. Create the worker script `src/index.ts`:

```typescript
import { Container } from "@cloudflare/containers";

interface Env {
  SMTP_BACKEND: any;
  SMTP_HOST: string;
  SMTP_PORT: string;
  SMTP_USERNAME: string;
  SMTP_PASSWORD: string;
  RECIPIENT_EMAIL: string;
  API_KEY: string;
}

class SMTPBackend extends Container {
  defaultPort = 8080;
  sleepAfter = "2h";
  
  envVars = {
    SMTP_HOST: env.SMTP_HOST,
    SMTP_PORT: env.SMTP_PORT,
    SMTP_USERNAME: env.SMTP_USERNAME,
    SMTP_PASSWORD: env.SMTP_PASSWORD,
    RECIPIENT_EMAIL: env.RECIPIENT_EMAIL,
    API_KEY: env.API_KEY,
    SERVER_PORT: "8080"
  };
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const containerInstance = env.SMTP_BACKEND.get(0);
    return containerInstance.fetch(request);
  },
};

export { SMTPBackend };
```

3. Set your secrets using Wrangler:

```bash
wrangler secret put SMTP_HOST
wrangler secret put SMTP_PORT
wrangler secret put SMTP_USERNAME
wrangler secret put SMTP_PASSWORD
wrangler secret put RECIPIENT_EMAIL
wrangler secret put API_KEY
```

4. Deploy to Cloudflare:

```bash
wrangler deploy
```

## Security Notes

- Always use environment variables for sensitive data
- The API key should be a long, random string
- Use app-specific passwords for email providers that support them
- Consider rate limiting in production environments