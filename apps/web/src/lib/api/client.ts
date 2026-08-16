import createClient from "openapi-fetch";

import type { paths } from "./schema.gen";

/**
 * Base URL of the Go API. Server Components talk to the API directly, so this
 * is a server-side variable (no NEXT_PUBLIC_ prefix).
 */
const baseUrl = process.env.API_BASE_URL ?? "http://localhost:8080";

/** Typed API client; every path, param and body comes from openapi/openapi.yaml. */
export const apiClient = createClient<paths>({ baseUrl });
