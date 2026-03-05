import { Outlet, Link } from '@tanstack/react-router'

export function RootLayout() {
  return (
    <div className="app">
      <header className="app-header">
        <h1>Travelist</h1>
        <p>旅行プランを簡単に作成</p>
        <nav className="app-nav">
          <Link to="/">ホーム</Link>
          <Link to="/register">ユーザー登録</Link>
        </nav>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}
