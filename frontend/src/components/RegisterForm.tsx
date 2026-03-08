import { useForm } from 'react-hook-form'
import { useMutation } from '@tanstack/react-query'
import { registerUser } from '../api/client'
import { ApiError } from '../api/errors'

interface FormValues {
  email: string
  password: string
  passwordConfirm: string
}

const membershipHighlights = [
  {
    title: '保存して再編集',
    description: '気になった旅程をアカウントにひも付けて、あとから微調整できます。',
  },
  {
    title: '更新を見逃さない',
    description: '機能追加や新しいプランニング体験をメールで受け取れます。',
  },
  {
    title: '次回の旅が速い',
    description: '好みの条件をもとに、次の旅行計画をスムーズに始められます。',
  },
]

const onboardingSteps = [
  'メールアドレスとパスワードを登録',
  '確認メールのリンクをクリック',
  'Travelist で次の旅程づくりを開始',
]

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

  return (
    <section className="register-shell">
      <div className="register-layout">
        <aside className="register-hero">
          <p className="register-eyebrow">Travelist Membership</p>
          <h2>旅のアイデアを、次の行動につなげるアカウント。</h2>
          <p className="register-lead">
            登録しておくと、確認メール経由で安全に利用を始められます。
            今後のアップデートや、繰り返し使う旅程作成もスムーズになります。
          </p>

          <div className="register-value-grid">
            {membershipHighlights.map((item) => (
              <article key={item.title} className="register-value-card">
                <h3>{item.title}</h3>
                <p>{item.description}</p>
              </article>
            ))}
          </div>

          <div className="register-steps-card">
            <p className="register-steps-title">登録の流れ</p>
            <ol className="register-steps">
              {onboardingSteps.map((step) => (
                <li key={step}>{step}</li>
              ))}
            </ol>
          </div>
        </aside>

        <div className="register-card">
          {mutation.isSuccess ? (
            <div className="register-success">
              <span className="register-status-badge">Ready</span>
              <h3>登録完了</h3>
              <p className="register-success-copy">
                確認メールを送信しました。メール内のリンクをクリックして登録を完了してください。
              </p>
              <div className="register-success-panel">
                <p className="register-success-label">送信先</p>
                <p className="register-success-email">{mutation.data.email}</p>
              </div>
              <p className="register-success-note">
                しばらく届かない場合は迷惑メールフォルダも確認してください。
              </p>
            </div>
          ) : (
            <>
              <div className="register-card-header">
                <p className="register-card-eyebrow">Create Account</p>
                <h3>ユーザー登録</h3>
                <p>
                  メール認証で利用を開始します。登録後に届く確認メールからアカウントを有効化してください。
                </p>
              </div>

              <form onSubmit={handleSubmit(onSubmit)} className="register-form">
                <div className="form-group">
                  <label htmlFor="email">メールアドレス</label>
                  <input
                    id="email"
                    type="email"
                    autoComplete="email"
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
                    autoComplete="new-password"
                    {...register('password', {
                      required: 'パスワードを入力してください',
                      minLength: { value: 8, message: 'パスワードは8文字以上で入力してください' },
                      maxLength: { value: 72, message: 'パスワードは72文字以下で入力してください' },
                    })}
                  />
                  <p className="field-hint">8文字以上72文字以下で設定してください。</p>
                  {errors.password && <span className="error">{errors.password.message}</span>}
                </div>

                <div className="form-group">
                  <label htmlFor="passwordConfirm">パスワード（確認）</label>
                  <input
                    id="passwordConfirm"
                    type="password"
                    autoComplete="new-password"
                    {...register('passwordConfirm', { required: 'パスワードを再入力してください' })}
                  />
                  {errors.passwordConfirm && <span className="error">{errors.passwordConfirm.message}</span>}
                </div>

                {errors.root && <div className="error register-form-error" role="alert">{errors.root.message}</div>}

                <button type="submit" className="register-submit" disabled={mutation.isPending}>
                  {mutation.isPending ? '登録中...' : '登録してはじめる'}
                </button>
              </form>
            </>
          )}
        </div>
      </div>
    </section>
  )
}
