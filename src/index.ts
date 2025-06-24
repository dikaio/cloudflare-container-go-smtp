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
  sleepAfter = "30s";
  
  constructor(state: DurableObjectState, env: Env) {
    super(state, env);
    this.envVars = {
      SMTP_HOST: env.SMTP_HOST,
      SMTP_PORT: env.SMTP_PORT,
      SMTP_USERNAME: env.SMTP_USERNAME,
      SMTP_PASSWORD: env.SMTP_PASSWORD,
      RECIPIENT_EMAIL: env.RECIPIENT_EMAIL,
      API_KEY: env.API_KEY,
      SERVER_PORT: "8080"
    };
  }
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const containerInstance = env.SMTP_BACKEND.get(0);
    return containerInstance.fetch(request);
  },
};

export { SMTPBackend };