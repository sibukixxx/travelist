import { Outlet, Link, useRouterState } from '@tanstack/react-router'

export function RootLayout() {
  const routerState = useRouterState()
  const currentPath = routerState.location.pathname

  return (
    <div className="app">
      <header className="app-header">
        <div className="app-header-inner">
          <div style={{ display: 'flex', alignItems: 'center' }}>
            <Link to="/" className="app-logo">
              <span className="app-logo-icon" aria-hidden="true">
                &#x2708;
              </span>
              <span className="app-logo-text">Travelist</span>
            </Link>
            <span className="app-tagline">
              旅のワクワクを、かたちにする
            </span>
          </div>
          <nav className="app-nav">
            <Link to="/" data-status={currentPath === '/' ? 'active' : undefined}>
              プランを作る
            </Link>
            <Link
              to="/register"
              data-status={currentPath === '/register' ? 'active' : undefined}
            >
              ユーザー登録
            </Link>
          </nav>
        </div>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
      <footer className="app-footer">
        Travelist &mdash; あなただけの旅プランを
      </footer>
    </div>
  )
}
