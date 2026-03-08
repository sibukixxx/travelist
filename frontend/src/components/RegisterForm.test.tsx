import { beforeEach, describe, expect, it, vi } from 'vitest'
import { screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { RegisterForm } from './RegisterForm'
import { renderWithProviders } from '../test/helpers'
import { ApiError } from '../api/errors'

vi.mock('../api/client', () => ({
  registerUser: vi.fn(),
}))

import { registerUser } from '../api/client'

const mockedRegisterUser = vi.mocked(registerUser)

describe('RegisterForm', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the onboarding content and form fields', () => {
    renderWithProviders(<RegisterForm />)

    expect(
      screen.getByRole('heading', { name: '旅のアイデアを、次の行動につなげるアカウント。' }),
    ).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'ユーザー登録' })).toBeInTheDocument()
    expect(screen.getByLabelText('メールアドレス')).toBeInTheDocument()
    expect(screen.getByLabelText('パスワード')).toBeInTheDocument()
    expect(screen.getByLabelText('パスワード（確認）')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '登録してはじめる' })).toBeInTheDocument()
  })

  it('shows validation error when password confirmation does not match', async () => {
    const user = userEvent.setup()
    renderWithProviders(<RegisterForm />)

    await user.type(screen.getByLabelText('メールアドレス'), 'user@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password123')
    await user.type(screen.getByLabelText('パスワード（確認）'), 'password321')
    await user.click(screen.getByRole('button', { name: '登録してはじめる' }))

    expect(await screen.findByText('パスワードが一致しません')).toBeInTheDocument()
    expect(mockedRegisterUser).not.toHaveBeenCalled()
  })

  it('shows API error inside the register card', async () => {
    const user = userEvent.setup()
    mockedRegisterUser.mockRejectedValueOnce(
      new ApiError(409, 'CONFLICT', 'メールアドレスは既に登録されています'),
    )
    renderWithProviders(<RegisterForm />)

    await user.type(screen.getByLabelText('メールアドレス'), 'user@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password123')
    await user.type(screen.getByLabelText('パスワード（確認）'), 'password123')
    await user.click(screen.getByRole('button', { name: '登録してはじめる' }))

    expect(await screen.findByRole('alert')).toHaveTextContent('このメールアドレスは既に登録されています')
  })

  it('shows the success card with the registered email', async () => {
    const user = userEvent.setup()
    mockedRegisterUser.mockResolvedValueOnce({
      user_id: 'user_1',
      email: 'user@example.com',
    })
    renderWithProviders(<RegisterForm />)

    await user.type(screen.getByLabelText('メールアドレス'), 'user@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password123')
    await user.type(screen.getByLabelText('パスワード（確認）'), 'password123')
    await user.click(screen.getByRole('button', { name: '登録してはじめる' }))

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: '登録完了' })).toBeInTheDocument()
    })
    expect(screen.getByText('user@example.com')).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '登録してはじめる' })).not.toBeInTheDocument()
  })
})
