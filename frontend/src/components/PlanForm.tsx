import { useForm, Controller } from 'react-hook-form'
import { useMutation } from '@tanstack/react-query'
import { generatePlan } from '../api/client'
import type { PlanRequest, GenerateResult } from '../types/itinerary'
import { TagInput } from './TagInput'

interface PlanFormProps {
  onResult: (result: GenerateResult) => void
}

interface FormValues {
  destination: string
  num_days: number
  start_date: string
  interests: string[]
  budget: 'budget' | 'moderate' | 'luxury'
  travel_style: 'relaxed' | 'active' | 'balanced'
  total_budget_yen: string
}

const interestSuggestions = [
  '文化', '歴史', '食事', '自然', 'アート', 'ショッピング',
  '温泉', '神社仏閣', '写真', 'アウトドア', 'グルメ', '夜景', '建築',
]

const today = new Date().toISOString().split('T')[0]

export function PlanForm({ onResult }: PlanFormProps) {
  const { register, handleSubmit, control, formState: { errors } } = useForm<FormValues>({
    defaultValues: {
      destination: '',
      num_days: 3,
      start_date: '',
      interests: [],
      budget: 'moderate',
      travel_style: 'balanced',
      total_budget_yen: '',
    },
  })

  const mutation = useMutation({
    mutationFn: generatePlan,
    onSuccess: onResult,
  })

  const onSubmit = (data: FormValues) => {
    const budgetNum = parseInt(data.total_budget_yen, 10)
    const req: PlanRequest = {
      destination: data.destination,
      num_days: data.num_days,
      start_date: data.start_date,
      preferences: {
        interests: data.interests,
        budget: data.budget,
        travel_style: data.travel_style,
        ...(Number.isFinite(budgetNum) && budgetNum > 0 ? { total_budget_yen: budgetNum } : {}),
      },
      constraint: {
        max_walk_distance_m: 2000,
        max_activities_day: 6,
        earliest_start_time: '08:00',
        latest_end_time: '21:00',
      },
    }
    mutation.mutate(req)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="plan-form">
      {/* Section 1: Destination */}
      <div className="form-section">
        <div className="form-section-header">
          <span className="form-section-icon form-section-icon--destination" aria-hidden="true">
            &#x1F4CD;
          </span>
          <h3 className="form-section-title">どこへ行く？</h3>
        </div>
        <div className="form-section-body">
          <div className="form-group">
            <label htmlFor="destination">目的地</label>
            <input
              id="destination"
              {...register('destination', { required: '目的地を入力してください' })}
              placeholder="例: 京都、沖縄、北海道..."
            />
            {errors.destination && <span className="error">{errors.destination.message}</span>}
          </div>
        </div>
      </div>

      {/* Section 2: Schedule */}
      <div className="form-section">
        <div className="form-section-header">
          <span className="form-section-icon form-section-icon--schedule" aria-hidden="true">
            &#x1F4C5;
          </span>
          <h3 className="form-section-title">いつ、何日間？</h3>
        </div>
        <div className="form-section-body">
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="start_date">開始日</label>
              <input
                id="start_date"
                type="date"
                min={today}
                {...register('start_date', {
                  required: '開始日を入力してください',
                  validate: (v) => v >= today || '過去の日付は選択できません',
                })}
              />
              {errors.start_date && <span className="error">{errors.start_date.message}</span>}
            </div>
            <div className="form-group">
              <label htmlFor="num_days">日数</label>
              <input
                id="num_days"
                type="number"
                {...register('num_days', { required: true, min: 1, max: 14, valueAsNumber: true })}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Section 3: Budget */}
      <div className="form-section">
        <div className="form-section-header">
          <span className="form-section-icon form-section-icon--budget" aria-hidden="true">
            &#x1F4B0;
          </span>
          <h3 className="form-section-title">予算はどのくらい？</h3>
        </div>
        <div className="form-section-body">
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="budget">予算</label>
              <select id="budget" {...register('budget')}>
                <option value="budget">節約</option>
                <option value="moderate">普通</option>
                <option value="luxury">贅沢</option>
              </select>
            </div>
            <div className="form-group">
              <label htmlFor="total_budget_yen">予算上限（円、任意）</label>
              <input
                id="total_budget_yen"
                type="number"
                {...register('total_budget_yen')}
                placeholder="例: 50000"
                min="0"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Section 4: Interests */}
      <div className="form-section">
        <div className="form-section-header">
          <span className="form-section-icon form-section-icon--interest" aria-hidden="true">
            &#x2728;
          </span>
          <h3 className="form-section-title">何に興味がある？</h3>
        </div>
        <div className="form-section-body">
          <div className="form-group">
            <label>興味・関心</label>
            <Controller
              name="interests"
              control={control}
              render={({ field }) => (
                <TagInput
                  tags={field.value}
                  onChange={field.onChange}
                  suggestions={interestSuggestions}
                />
              )}
            />
          </div>
        </div>
      </div>

      {/* Section 5: Travel Style */}
      <div className="form-section">
        <div className="form-section-header">
          <span className="form-section-icon form-section-icon--style" aria-hidden="true">
            &#x1F6B6;
          </span>
          <h3 className="form-section-title">どんなペースで？</h3>
        </div>
        <div className="form-section-body">
          <div className="form-group">
            <label htmlFor="travel_style">旅行スタイル</label>
            <select id="travel_style" {...register('travel_style')}>
              <option value="relaxed">ゆったり</option>
              <option value="balanced">バランス</option>
              <option value="active">アクティブ</option>
            </select>
          </div>
        </div>
      </div>

      {/* Submit */}
      <div className="form-submit">
        <button type="submit" disabled={mutation.isPending}>
          {mutation.isPending ? 'プラン生成中...' : 'プランを生成'}
        </button>
        <p className="form-submit-hint">
          AI が目的地・日程・予算に合わせた旅行プランを自動作成します
        </p>
      </div>

      {mutation.isError && (
        <div className="error">エラー: {mutation.error.message}</div>
      )}
    </form>
  )
}
