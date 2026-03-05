export interface RegisterRequest {
  email: string
  password: string
}

export interface RegisterResponse {
  user_id: string
  email: string
}
