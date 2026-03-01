import { Outlet } from '@tanstack/react-router'

export function RootLayout() {
  return (
    <div className="app">
      <header className="app-header">
        <h1>Travelist</h1>
        <p>旅行プランを簡単に作成</p>
      </header>
      <main className="app-main">
        <Outlet />
      </main>
    </div>
  )
}
