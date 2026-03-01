import { useForm } from 'react-hook-form'
import { useMutation } from '@tanstack/react-query'
import { registerUser } from '../api/client'
import { ApiError } from '../api/errors'

interface FormValues {
  email: string
  password: string
  passwordConfirm: string
}

export function RegisterForm() {
  const { register, handleSubmit, formState: { errors }, setError } = useForm<FormValues>()

  const mutation = useMutation({
    mutationFn: registerUser,
    onError: (error) => {
      if (error instanceof ApiError) {
        setError('root', { message: error.message })
      } else {
        setError('root', { message: 'エラーが発生しました' })
      }
    },
  })

  const onSubmit = (data: FormValues) => {
    if (data.password !== data.passwordConfirm) {
      setError('passwordConfirm', { message: 'パスワードが一致しません' })
      return
    }
    mutation.mutate({ email: data.email, password: data.password })
  }

  if (mutation.isSuccess) {
    return (
      <div className="register-success">
        <h2>登録完了</h2>
        <p>確認メールを送信しました。メール内のリンクをクリックして登録を完了してください。</p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="register-form">
      <h2>ユーザー登録</h2>

      <div className="form-group">
        <label htmlFor="email">メールアドレス</label>
        <input
          id="email"
          type="email"
          {...register('email', { required: 'メールアドレスを入力してください' })}
          placeholder="example@mail.com"
        />
        {errors.email && <span className="error">{errors.email.message}</span>}
      </div>

      <div className="form-group">
        <label htmlFor="password">パスワード</label>
        <input
          id="password"
          type="password"
          {...register('password', {
            required: 'パスワードを入力してください',
            minLength: { value: 8, message: 'パスワードは8文字以上で入力してください' },
            maxLength: { value: 72, message: 'パスワードは72文字以下で入力してください' },
          })}
        />
        {errors.password && <span className="error">{errors.password.message}</span>}
      </div>

      <div className="form-group">
        <label htmlFor="passwordConfirm">パスワード（確認）</label>
        <input
          id="passwordConfirm"
          type="password"
          {...register('passwordConfirm', { required: 'パスワードを再入力してください' })}
        />
        {errors.passwordConfirm && <span className="error">{errors.passwordConfirm.message}</span>}
      </div>

      {errors.root && <div className="error">{errors.root.message}</div>}

      <button type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? '登録中...' : '登録'}
      </button>
    </form>
  )
}
