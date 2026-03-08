import { useRef, useState } from 'react'
import { PlanForm } from './PlanForm'
import { BudgetSummaryDisplay } from './BudgetSummary'
import { ResultActions } from './ResultActions'
import type { GenerateResult } from '../types/itinerary'

export function HomePage() {
  const [result, setResult] = useState<GenerateResult | null>(null)
  const formRef = useRef<HTMLDivElement>(null)
  const resultRef = useRef<HTMLDivElement>(null)

  const handleResult = (data: GenerateResult) => {
    setResult(data)
    setTimeout(() => {
      resultRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, 100)
  }

  const handleRegenerate = () => {
    formRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div>
      {!result && (
        <div className="hero">
          <h1 className="hero-title">旅の計画を、もっと楽しく</h1>
          <p className="hero-subtitle">
            行き先と日程を入れるだけで、あなたにぴったりの旅プランをご提案します
          </p>
        </div>
      )}

      <div ref={formRef} className="card card-accent" style={{ marginBottom: '2rem' }}>
        <h2 className="section-title">プラン条件</h2>
        <PlanForm onResult={handleResult} />
      </div>

      {result && (
        <div ref={resultRef} className="result">
          <div className="result-header">
            <h2 className="result-title">{result.itinerary.title}</h2>
            <ResultActions result={result} onRegenerate={handleRegenerate} />
          </div>

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
                <h3>Day {day.day_number} &mdash; {day.date}</h3>
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
