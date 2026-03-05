import type { GenerateResult } from '../types/itinerary'
import { downloadJson } from '../utils/downloadJson'
import { buildItineraryFilename } from '../utils/buildItineraryFilename'

interface ResultActionsProps {
  result: GenerateResult
  onRegenerate: () => void
}

export function ResultActions({ result, onRegenerate }: ResultActionsProps) {
  const handleDownload = () => {
    const filename = buildItineraryFilename(
      result.itinerary.destination,
      result.itinerary.start_date,
    )
    downloadJson(result, filename)
  }

  return (
    <div className="result-actions">
      <button type="button" className="btn-secondary" onClick={handleDownload}>
        JSONダウンロード
      </button>
      <button type="button" className="btn-secondary" onClick={onRegenerate}>
        プランを再生成
      </button>
    </div>
  )
}
