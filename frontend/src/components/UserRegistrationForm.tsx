import { useForm } from 'react-hook-form'
import { useMutation } from '@tanstack/react-query'
import { registerUser } from '../api/client'

interface FormValues {
  email: string
}

export function UserRegistrationForm() {
  const { register, handleSubmit, reset, formState: { errors } } = useForm<FormValues>({
    defaultValues: { email: '' },
  })

  const mutation = useMutation({
    mutationFn: (email: string) => registerUser(email),
    onSuccess: () => {
      reset()
    },
  })

  const onSubmit = (data: FormValues) => {
    mutation.mutate(data.email)
  }

  return (
    <section className="registration-card">
      <h2>ユーザ登録</h2>
      <p className="registration-help">メールアドレスを登録してアップデート通知を受け取れます。</p>
      <form onSubmit={handleSubmit(onSubmit)} className="registration-form">
        <div className="form-group">
          <label htmlFor="email">メールアドレス</label>
          <input
            id="email"
            type="email"
            autoComplete="email"
            placeholder="you@example.com"
            {...register('email', {
              required: 'メールアドレスを入力してください',
              pattern: {
                value: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                message: 'メールアドレスの形式が不正です',
              },
            })}
          />
          {errors.email && <span className="error">{errors.email.message}</span>}
        </div>

        <button type="submit" disabled={mutation.isPending}>
          {mutation.isPending ? '登録中...' : 'メール登録'}
        </button>
      </form>

      {mutation.isSuccess && (
        <p className="success-message">
          {mutation.data.email} を登録しました。
        </p>
      )}
      {mutation.isError && (
        <p className="error">登録に失敗しました: {mutation.error.message}</p>
      )}
    </section>
  )
}
