import { useState } from 'react'
import { PlanForm } from './PlanForm'
import type { GenerateResult } from '../types/itinerary'

export function HomePage() {
  const [result, setResult] = useState<GenerateResult | null>(null)

  return (
    <div>
      <PlanForm onResult={setResult} />
      {result && (
        <div className="result">
          <h2>{result.itinerary.title}</h2>
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
          {result.itinerary.days.map((day) => (
            <div key={day.day_number} className="day-plan">
              <h3>Day {day.day_number} - {day.date}</h3>
              <ul>
                {day.activities.map((act) => (
                  <li key={act.order}>
                    <strong>{act.start_time}–{act.end_time}</strong>{' '}
                    {act.place?.name ?? act.place_id}
                    {act.note && <span className="note"> — {act.note}</span>}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
