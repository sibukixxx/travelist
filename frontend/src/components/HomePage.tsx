import { useRef, useState } from 'react'
import { PlanForm } from './PlanForm'
import { BudgetSummaryDisplay } from './BudgetSummary'
import { ResultActions } from './ResultActions'
import type { GenerateResult } from '../types/itinerary'

export function HomePage() {
  const [result, setResult] = useState<GenerateResult | null>(null)
  const formRef = useRef<HTMLDivElement>(null)

  const handleRegenerate = () => {
    formRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div>
      <div ref={formRef}>
        <PlanForm onResult={setResult} />
      </div>
      {result && (
        <div className="result">
          <h2>{result.itinerary.title}</h2>
          <ResultActions result={result} onRegenerate={handleRegenerate} />
          {result.violations.length > 0 && (
            <div className="violations">
              <h3>注意事項</h3>
              <ul>
                {result.violations.map((v, i) => (
                  <li key={i}>{v.message}</li>
                ))}
              </ul>
            </div>
          )}
          {result.budget_summary && (
            <BudgetSummaryDisplay
              summary={result.budget_summary}
              violations={result.violations}
            />
          )}
          {result.itinerary.days.map((day) => {
            const dayCost = result.budget_summary?.daily_costs.find(
              (dc) => dc.day_number === day.day_number
            )
            return (
              <div key={day.day_number} className="day-plan">
                <h3>Day {day.day_number} - {day.date}</h3>
                <ul>
                  {day.activities.map((act) => (
                    <li key={act.order}>
                      <strong>{act.start_time}–{act.end_time}</strong>{' '}
                      {act.place?.name ?? act.place_id}
                      {act.estimated_cost_yen > 0 && (
                        <span className="cost"> ({act.estimated_cost_yen.toLocaleString('ja-JP')}円)</span>
                      )}
                      {act.note && <span className="note"> — {act.note}</span>}
                    </li>
                  ))}
                </ul>
                {dayCost && dayCost.cost_yen > 0 && (
                  <p className="day-cost-subtotal">
                    小計: {dayCost.cost_yen.toLocaleString('ja-JP')}円
                  </p>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
