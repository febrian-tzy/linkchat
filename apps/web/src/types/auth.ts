export interface RegisterPayload {
  workspace_name: string;
  full_name: string;
  email: string;
  username: string;
  password: string;
  confirm_password: string;
}

export interface LoginPayload {
  username: string;
  password: string;
}

export interface MeResponse {
  user_id: number;
  tenant_id: number;
  full_name: string;
  username: string;
  email: string;
  role: string;
  tenant: {
    id: number;
    name: string;
    slug: string;
  };
}

export interface ApiError {
  error: string;
}
