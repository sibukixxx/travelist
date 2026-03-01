import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { BudgetSummaryDisplay } from './BudgetSummary'
import type { BudgetSummary, Violation } from '../types/itinerary'

describe('BudgetSummaryDisplay', () => {
  it('renders formatted total cost', () => {
    const summary: BudgetSummary = {
      total_cost_yen: 5500,
      daily_costs: [{ day_number: 1, cost_yen: 5500 }],
    }

    render(<BudgetSummaryDisplay summary={summary} violations={[]} />)

    expect(screen.getByText('5,500円')).toBeInTheDocument()
  })

  it('shows warning when budget exceeded violation exists', () => {
    const summary: BudgetSummary = {
      total_cost_yen: 15000,
      daily_costs: [{ day_number: 1, cost_yen: 15000 }],
    }
    const violations: Violation[] = [
      {
        type: 'budget_exceeded',
        day_number: 0,
        activity_idx: 0,
        message: 'total estimated cost 15000円 exceeds budget 10000円',
      },
    ]

    render(<BudgetSummaryDisplay summary={summary} violations={violations} />)

    expect(screen.getByText(/予算超過/)).toBeInTheDocument()
  })

  it('does not show warning when no budget violation', () => {
    const summary: BudgetSummary = {
      total_cost_yen: 5000,
      daily_costs: [{ day_number: 1, cost_yen: 5000 }],
    }

    render(<BudgetSummaryDisplay summary={summary} violations={[]} />)

    expect(screen.queryByText(/予算超過/)).not.toBeInTheDocument()
  })
})
