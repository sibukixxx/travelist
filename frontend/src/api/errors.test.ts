import { describe, it, expect } from 'vitest'
import { ApiError, parseApiError } from './errors'

describe('ApiError', () => {
  it('stores status, code, and serverMessage', () => {
    const err = new ApiError(409, 'CONFLICT', 'email already exists')
    expect(err.status).toBe(409)
    expect(err.code).toBe('CONFLICT')
    expect(err.serverMessage).toBe('email already exists')
    expect(err).toBeInstanceOf(Error)
  })

  it('provides Japanese message for CONFLICT', () => {
    const err = new ApiError(409, 'CONFLICT', 'email already exists')
    expect(err.message).toBe('このメールアドレスは既に登録されています')
  })

  it('provides Japanese message for BAD_REQUEST', () => {
    const err = new ApiError(400, 'BAD_REQUEST', 'bad input')
    expect(err.message).toBe('入力内容に誤りがあります')
  })

  it('provides Japanese message for VALIDATION_ERROR', () => {
    const err = new ApiError(400, 'VALIDATION_ERROR', 'field error')
    expect(err.message).toBe('入力値が正しくありません')
  })

  it('provides Japanese message for INVALID_TOKEN', () => {
    const err = new ApiError(400, 'INVALID_TOKEN', 'token is invalid')
    expect(err.message).toBe('認証トークンが無効です')
  })

  it('provides Japanese message for EXPIRED_TOKEN', () => {
    const err = new ApiError(400, 'EXPIRED_TOKEN', 'token expired')
    expect(err.message).toBe('認証トークンの有効期限が切れています')
  })

  it('provides Japanese message for INTERNAL', () => {
    const err = new ApiError(500, 'INTERNAL', 'internal server error')
    expect(err.message).toBe('サーバーエラーが発生しました。しばらくしてからお試しください')
  })

  it('falls back to server message for unknown code', () => {
    const err = new ApiError(418, 'TEAPOT', 'I am a teapot')
    expect(err.message).toBe('I am a teapot')
  })
})

describe('parseApiError', () => {
  it('parses valid JSON error response', () => {
    const body = JSON.stringify({ status: 409, code: 'CONFLICT', message: 'duplicate' })
    const err = parseApiError(409, body)
    expect(err).toBeInstanceOf(ApiError)
    expect(err.status).toBe(409)
    expect(err.code).toBe('CONFLICT')
    expect(err.serverMessage).toBe('duplicate')
  })

  it('returns ApiError with BAD_REQUEST for invalid JSON', () => {
    const err = parseApiError(400, 'not json')
    expect(err).toBeInstanceOf(ApiError)
    expect(err.status).toBe(400)
    expect(err.code).toBe('UNKNOWN')
    expect(err.serverMessage).toBe('not json')
  })

  it('handles JSON without code field', () => {
    const body = JSON.stringify({ error: 'something went wrong' })
    const err = parseApiError(500, body)
    expect(err).toBeInstanceOf(ApiError)
    expect(err.status).toBe(500)
    expect(err.code).toBe('UNKNOWN')
  })
})
