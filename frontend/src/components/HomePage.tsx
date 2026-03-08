import { useRef, useState } from 'react'
import { PlanForm } from './PlanForm'
import { BudgetSummaryDisplay } from './BudgetSummary'
import { ResultActions } from './ResultActions'
import { UserRegistrationForm } from './UserRegistrationForm'
import type { GenerateResult } from '../types/itinerary'

const travelModeLabels: Record<string, string> = {
  walk: '徒歩',
  train: '電車',
  bus: 'バス',
  taxi: 'タクシー',
  driving: '車',
}

const travelModeIcons: Record<string, string> = {
  walk: '\u{1F6B6}',
  train: '\u{1F683}',
  bus: '\u{1F68C}',
  taxi: '\u{1F695}',
  driving: '\u{1F697}',
}

export function HomePage() {
  const [result, setResult] = useState<GenerateResult | null>(null)
  const formRef = useRef<HTMLDivElement>(null)

  const handleRegenerate = () => {
    formRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  return (
    <div>
      {/* Hero */}
      <section className="hero">
        <div className="hero-inner">
          <h2>旅の計画を、もっと楽しく</h2>
          <p>行き先と日程を入力するだけで、あなただけの旅行プランを提案します。</p>
          <div className="hero-stats">
            <div className="hero-stat">
              <span className="hero-stat-value">100+</span>
              <span className="hero-stat-label">観光スポット</span>
            </div>
            <div className="hero-stat">
              <span className="hero-stat-value">3min</span>
              <span className="hero-stat-label">プラン生成</span>
            </div>
            <div className="hero-stat">
              <span className="hero-stat-value">Free</span>
              <span className="hero-stat-label">利用料金</span>
            </div>
          </div>
        </div>
      </section>

      {/* Builder */}
      <div className="builder-section" ref={formRef}>
        <PlanForm onResult={setResult} />
      </div>

      {/* Results */}
      {result && (
        <div className="result">
          <div className="result-hero">
            <h2>{result.itinerary.title}</h2>
            <div className="result-meta">
              <span>{result.itinerary.destination}</span>
              <span>{result.itinerary.start_date} - {result.itinerary.end_date}</span>
              <span>{result.itinerary.days.length}日間</span>
            </div>
          </div>

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
                <div className="day-header">
                  <div className="day-badge">{day.day_number}</div>
                  <div className="day-header-text">
                    <h3>Day {day.day_number}</h3>
                    <span className="day-date">{day.date}</span>
                  </div>
                </div>
                <div className="timeline">
                  {day.activities.map((act) => (
                    <div key={act.order} className="timeline-item">
                      <span className="timeline-time">
                        {act.start_time} - {act.end_time}
                      </span>
                      <div>
                        <span className="timeline-place">
                          {act.place?.name ?? act.place_id}
                        </span>
                        {act.estimated_cost_yen > 0 && (
                          <span className="timeline-cost">
                            {act.estimated_cost_yen.toLocaleString('ja-JP')}円
                          </span>
                        )}
                      </div>
                      {act.note && (
                        <span className="timeline-note">{act.note}</span>
                      )}
                      {act.travel_from_prev && (
                        <div className="timeline-travel">
                          <span className="timeline-travel-icon">
                            {travelModeIcons[act.travel_from_prev.mode] ?? ''}
                          </span>
                          <span>
                            {travelModeLabels[act.travel_from_prev.mode] ?? act.travel_from_prev.mode}
                            {' '}{act.travel_from_prev.duration_min}分
                          </span>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
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

      {/* Registration CTA */}
      <UserRegistrationForm />
    </div>
  )
}
