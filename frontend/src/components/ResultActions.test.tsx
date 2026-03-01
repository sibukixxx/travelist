import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ResultActions } from './ResultActions'
import type { GenerateResult } from '../types/itinerary'
import * as downloadJsonModule from '../utils/downloadJson'

vi.mock('../utils/downloadJson', () => ({
  downloadJson: vi.fn(),
}))

const mockResult: GenerateResult = {
  itinerary: {
    id: '1',
    title: 'Tokyo Trip',
    destination: 'Tokyo',
    start_date: '2025-07-01',
    end_date: '2025-07-03',
    days: [],
    created_at: '',
    updated_at: '',
  },
  violations: [],
  budget_summary: { total_cost_yen: 10000, daily_costs: [] },
}

describe('ResultActions', () => {
  it('renders download and regenerate buttons', () => {
    render(<ResultActions result={mockResult} onRegenerate={vi.fn()} />)

    expect(screen.getByRole('button', { name: 'JSONダウンロード' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'プランを再生成' })).toBeInTheDocument()
  })

  it('calls downloadJson with result data when download button is clicked', async () => {
    const user = userEvent.setup()
    render(<ResultActions result={mockResult} onRegenerate={vi.fn()} />)

    await user.click(screen.getByRole('button', { name: 'JSONダウンロード' }))

    expect(downloadJsonModule.downloadJson).toHaveBeenCalledWith(
      mockResult,
      'travelist-Tokyo-2025-07-01.json',
    )
  })

  it('calls onRegenerate when regenerate button is clicked', async () => {
    const user = userEvent.setup()
    const onRegenerate = vi.fn()
    render(<ResultActions result={mockResult} onRegenerate={onRegenerate} />)

    await user.click(screen.getByRole('button', { name: 'プランを再生成' }))

    expect(onRegenerate).toHaveBeenCalledOnce()
  })
})
