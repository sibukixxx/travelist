const ERROR_MESSAGES: Record<string, string> = {
  CONFLICT: 'このメールアドレスは既に登録されています',
  BAD_REQUEST: '入力内容に誤りがあります',
  VALIDATION_ERROR: '入力値が正しくありません',
  INVALID_TOKEN: '認証トークンが無効です',
  EXPIRED_TOKEN: '認証トークンの有効期限が切れています',
  INTERNAL: 'サーバーエラーが発生しました。しばらくしてからお試しください',
  NOT_FOUND: 'リソースが見つかりません',
}

export class ApiError extends Error {
  readonly status: number
  readonly code: string
  readonly serverMessage: string

  constructor(status: number, code: string, serverMessage: string) {
    const message = ERROR_MESSAGES[code] ?? serverMessage
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.serverMessage = serverMessage
  }
}

interface ErrorResponseBody {
  status?: number
  code?: string
  message?: string
}

export function parseApiError(status: number, body: string): ApiError {
  try {
    const parsed: ErrorResponseBody = JSON.parse(body)
    const code = parsed.code ?? 'UNKNOWN'
    const serverMessage = parsed.message ?? body
    return new ApiError(status, code, serverMessage)
  } catch {
    return new ApiError(status, 'UNKNOWN', body)
  }
}
