import type { LoginPayload, MeResponse, RegisterPayload } from "@/types/auth";

const API_URL = "http://localhost:8080";
const BASE = `${API_URL}/api/v1`;

const TOKEN_KEY = "linkchat_token";

export const tokenStorage = {
  get: (): string | null => {
    if (typeof window === "undefined") return null;
    return localStorage.getItem(TOKEN_KEY);
  },
  set: (token: string) => {
    if (typeof window === "undefined") return;
    localStorage.setItem(TOKEN_KEY, token);
  },
  clear: () => {
    if (typeof window === "undefined") return;
    localStorage.removeItem(TOKEN_KEY);
  },
};

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const token = tokenStorage.get();
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  let res: Response;
  try {
    res = await fetch(`${BASE}${path}`, { headers, ...options });
  } catch {
    throw new Error("Tidak dapat terhubung ke server. Pastikan backend sudah berjalan.");
  }

  let data: unknown;
  const ct = res.headers.get("content-type") ?? "";
  if (ct.includes("application/json")) {
    data = await res.json();
  } else {
    const text = await res.text();
    data = { error: text || "Server error" };
  }

  if (!res.ok) {
    throw new Error((data as Record<string, string>)?.error ?? "Something went wrong");
  }
  return data as T;
}

// ── Auth ──────────────────────────────────────────────────────
export const authApi = {
  register: (payload: RegisterPayload) =>
    request("/auth/register", { method: "POST", body: JSON.stringify(payload) }),

  login: async (payload: LoginPayload) => {
    const data = await request<{ message: string; token: string }>("/auth/login", {
      method: "POST", body: JSON.stringify(payload),
    });
    if (data.token) tokenStorage.set(data.token);
    return data;
  },

  logout: () => { tokenStorage.clear(); return Promise.resolve(); },
  me: () => request<MeResponse>("/auth/me"),
};

// ── Contacts ──────────────────────────────────────────────────
export const contactsApi = {
  list: (params?: { search?: string; limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.search) q.set("search", params.search);
    if (params?.limit)  q.set("limit",  String(params.limit));
    if (params?.offset) q.set("offset", String(params.offset));
    return request<{ data: Contact[]; total: number }>(`/contacts?${q}`);
  },
  get:    (id: number) => request<Contact>(`/contacts/${id}`),
  create: (body: Partial<Contact>) => request<Contact>("/contacts", { method: "POST", body: JSON.stringify(body) }),
  update: (id: number, body: Partial<Contact>) => request<Contact>(`/contacts/${id}`, { method: "PUT", body: JSON.stringify(body) }),
  delete: (id: number) => request(`/contacts/${id}`, { method: "DELETE" }),
};

// ── Inboxes ───────────────────────────────────────────────────
export const inboxesApi = {
  list:   () => request<{ data: Inbox[] }>("/inboxes"),
  get:    (id: number) => request<Inbox>(`/inboxes/${id}`),
  create: (body: { name: string; channel_type: string }) => request<Inbox>("/inboxes", { method: "POST", body: JSON.stringify(body) }),
  delete: (id: number) => request(`/inboxes/${id}`, { method: "DELETE" }),
};

// ── Conversations ─────────────────────────────────────────────
export const conversationsApi = {
  list: (params?: { status?: string; limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.status) q.set("status", params.status);
    if (params?.limit)  q.set("limit",  String(params.limit));
    if (params?.offset) q.set("offset", String(params.offset));
    return request<{ data: Conversation[]; total: number }>(`/conversations?${q}`);
  },
  get:          (id: number) => request<Conversation>(`/conversations/${id}`),
  create:       (body: { inbox_id: number; contact_id: number }) => request<Conversation>("/conversations", { method: "POST", body: JSON.stringify(body) }),
  updateStatus: (id: number, body: { status: string; assigned_user_id?: number }) =>
    request(`/conversations/${id}`, { method: "PATCH", body: JSON.stringify(body) }),
  listMessages: (id: number, params?: { limit?: number; offset?: number }) => {
    const q = new URLSearchParams();
    if (params?.limit)  q.set("limit",  String(params.limit));
    if (params?.offset) q.set("offset", String(params.offset));
    return request<{ data: Message[]; total: number }>(`/conversations/${id}/messages?${q}`);
  },
  sendMessage: (id: number, body: { content: string; content_type?: string }) =>
    request<Message>(`/conversations/${id}/messages`, { method: "POST", body: JSON.stringify(body) }),
};

// ── Broadcasts ────────────────────────────────────────────────
export const broadcastsApi = {
  list:   () => request<{ data: Broadcast[] }>("/broadcasts"),
  get:    (id: number) => request<Broadcast>(`/broadcasts/${id}`),
  create: (body: Partial<Broadcast> & { contact_ids: number[] }) =>
    request<Broadcast>("/broadcasts", { method: "POST", body: JSON.stringify(body) }),
  delete: (id: number) => request(`/broadcasts/${id}`, { method: "DELETE" }),
};

// ── Webhooks & API Keys ───────────────────────────────────────
export const webhooksApi = {
  list:   () => request<{ data: Webhook[] }>("/webhooks"),
  create: (body: { url: string; secret?: string; events: string[] }) =>
    request<Webhook>("/webhooks", { method: "POST", body: JSON.stringify(body) }),
  delete: (id: number) => request(`/webhooks/${id}`, { method: "DELETE" }),
};

export const apiKeysApi = {
  list:   () => request<{ data: APIKey[] }>("/api-keys"),
  create: (body: { name: string }) =>
    request<APIKey & { key: string }>("/api-keys", { method: "POST", body: JSON.stringify(body) }),
  delete: (id: number) => request(`/api-keys/${id}`, { method: "DELETE" }),
};

// ── Automations ───────────────────────────────────────────────
export const automationsApi = {
  list:   () => request<{ data: Automation[] }>("/automations"),
  get:    (id: number) => request<Automation>(`/automations/${id}`),
  create: (body: Partial<Automation>) =>
    request<Automation>("/automations", { method: "POST", body: JSON.stringify(body) }),
  toggle: (id: number, is_active: boolean) =>
    request(`/automations/${id}/toggle`, { method: "PATCH", body: JSON.stringify({ is_active }) }),
  delete: (id: number) => request(`/automations/${id}`, { method: "DELETE" }),
};

// ── Users ─────────────────────────────────────────────────────
export const usersApi = {
  list:    () => request<{ data: TeamUser[] }>("/users"),
  invite:  (body: { full_name: string; email: string; username: string; password: string; role: string }) =>
    request("/users/invite", { method: "POST", body: JSON.stringify(body) }),
  updateProfile: (body: { full_name: string; email: string }) =>
    request("/users/profile", { method: "PUT", body: JSON.stringify(body) }),
  changePassword: (body: { current_password: string; new_password: string }) =>
    request("/users/change-password", { method: "POST", body: JSON.stringify(body) }),
  toggle: (id: number, is_active: boolean) =>
    request(`/users/${id}/toggle`, { method: "PATCH", body: JSON.stringify({ is_active }) }),
  delete: (id: number) => request(`/users/${id}`, { method: "DELETE" }),
};

// ── Types ─────────────────────────────────────────────────────
export interface Contact {
  id: number; tenant_id: number; name: string; phone: string;
  email: string; notes: string; is_blocked: boolean;
  created_at: string; updated_at: string;
}
export interface Inbox {
  id: number; tenant_id: number; name: string; channel_type: string;
  is_enabled: boolean; created_at: string;
}
export interface Conversation {
  id: number; tenant_id: number; inbox_id: number; contact_id: number;
  status: string; unread_count: number; last_message_at: string;
  contact?: { name: string; phone: string };
  inbox?: { name: string };
  created_at: string;
}
export interface Message {
  id: number; conversation_id: number; sender_type: string;
  content: string; content_type: string; is_read: boolean; created_at: string;
}
export interface Broadcast {
  id: number; name: string; message: string; status: string;
  inbox_id: number; total: number; sent_count: number; failed_count: number;
  created_at: string;
}
export interface Webhook {
  id: number; url: string; events: string[]; is_active: boolean; created_at: string;
}
export interface APIKey {
  id: number; name: string; key_prefix: string; is_active: boolean; created_at: string;
}
export interface Automation {
  id: number; name: string; description: string; event: string;
  conditions: unknown; actions: unknown; is_active: boolean; run_count: number;
}
export interface TeamUser {
  id: number; full_name: string; email: string; username: string;
  role: string; is_active: boolean;
}

// ── Owner (Platform Admin) ────────────────────────────────────
const OWNER_TOKEN_KEY = "linkchat_owner_token";

export const ownerTokenStorage = {
  get:   (): string | null => typeof window !== "undefined" ? localStorage.getItem(OWNER_TOKEN_KEY) : null,
  set:   (t: string) => typeof window !== "undefined" && localStorage.setItem(OWNER_TOKEN_KEY, t),
  clear: () => typeof window !== "undefined" && localStorage.removeItem(OWNER_TOKEN_KEY),
};

async function ownerRequest<T>(path: string, options?: RequestInit): Promise<T> {
  const token = ownerTokenStorage.get();
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (token) headers["Authorization"] = `Bearer ${token}`;

  let res: Response;
  try {
    res = await fetch(`${API_URL}/api/v1/owner${path}`, { headers, ...options });
  } catch {
    throw new Error("Cannot connect to server.");
  }

  const ct = res.headers.get("content-type") ?? "";
  let data: unknown;
  if (ct.includes("application/json")) {
    data = await res.json();
  } else {
    data = { error: await res.text() || "Server error" };
  }

  if (!res.ok) throw new Error((data as Record<string, string>)?.error ?? "Error");
  return data as T;
}

export interface OwnerMe {
  id: number; full_name: string; email: string; username: string; role: string;
}
export interface PlatformStats {
  total_tenants: number; active_tenants: number;
  total_users: number; total_contacts: number; total_conversations: number;
}

export const ownerApi = {
  login: async (username: string, password: string) => {
    const data = await ownerRequest<{ token: string }>("/auth/login", {
      method: "POST", body: JSON.stringify({ username, password }),
    });
    if (data.token) ownerTokenStorage.set(data.token);
    return data;
  },
  logout: () => { ownerTokenStorage.clear(); return Promise.resolve(); },
  me:          () => ownerRequest<OwnerMe>("/auth/me"),
  stats:       () => ownerRequest<PlatformStats>("/stats"),
  listTenants: () => ownerRequest<{ data: Tenant[]; total: number }>("/tenants"),
};

export interface Tenant {
  id: number; name: string; slug: string; status: string; created_at: string;
}
