import type { BudgetSummary, Violation } from '../types/itinerary'

interface BudgetSummaryDisplayProps {
  summary: BudgetSummary
  violations: Violation[]
}

function formatYen(amount: number): string {
  return amount.toLocaleString('ja-JP') + '円'
}

export function BudgetSummaryDisplay({ summary, violations }: BudgetSummaryDisplayProps) {
  const budgetExceeded = violations.some((v) => v.type === 'budget_exceeded')

  return (
    <div className="budget-summary">
      <h3>予算サマリー</h3>
      <p className="total-cost">
        合計: <strong>{formatYen(summary.total_cost_yen)}</strong>
      </p>
      {budgetExceeded && (
        <p className="budget-over">予算超過: 設定した予算を超えています</p>
      )}
    </div>
  )
}
