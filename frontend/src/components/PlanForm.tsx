import { useForm } from 'react-hook-form'
import { useMutation } from '@tanstack/react-query'
import { generatePlan } from '../api/client'
import type { PlanRequest, GenerateResult } from '../types/itinerary'

interface PlanFormProps {
  onResult: (result: GenerateResult) => void
}

interface FormValues {
  destination: string
  num_days: number
  start_date: string
  interests: string
  budget: 'budget' | 'moderate' | 'luxury'
  travel_style: 'relaxed' | 'active' | 'balanced'
}

export function PlanForm({ onResult }: PlanFormProps) {
  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>({
    defaultValues: {
      destination: '',
      num_days: 3,
      start_date: '',
      interests: '',
      budget: 'moderate',
      travel_style: 'balanced',
    },
  })

  const mutation = useMutation({
    mutationFn: generatePlan,
    onSuccess: onResult,
  })

  const onSubmit = (data: FormValues) => {
    const req: PlanRequest = {
      destination: data.destination,
      num_days: data.num_days,
      start_date: data.start_date,
      preferences: {
        interests: data.interests.split(',').map((s) => s.trim()).filter(Boolean),
        budget: data.budget,
        travel_style: data.travel_style,
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
      <div className="form-group">
        <label htmlFor="destination">目的地</label>
        <input
          id="destination"
          {...register('destination', { required: '目的地を入力してください' })}
          placeholder="例: 京都"
        />
        {errors.destination && <span className="error">{errors.destination.message}</span>}
      </div>

      <div className="form-group">
        <label htmlFor="num_days">日数</label>
        <input
          id="num_days"
          type="number"
          {...register('num_days', { required: true, min: 1, max: 14, valueAsNumber: true })}
        />
      </div>

      <div className="form-group">
        <label htmlFor="start_date">開始日</label>
        <input
          id="start_date"
          type="date"
          {...register('start_date', { required: '開始日を入力してください' })}
        />
        {errors.start_date && <span className="error">{errors.start_date.message}</span>}
      </div>

      <div className="form-group">
        <label htmlFor="interests">興味・関心（カンマ区切り）</label>
        <input
          id="interests"
          {...register('interests')}
          placeholder="例: 文化, 食事, 自然"
        />
      </div>

      <div className="form-group">
        <label htmlFor="budget">予算</label>
        <select id="budget" {...register('budget')}>
          <option value="budget">節約</option>
          <option value="moderate">普通</option>
          <option value="luxury">贅沢</option>
        </select>
      </div>

      <div className="form-group">
        <label htmlFor="travel_style">旅行スタイル</label>
        <select id="travel_style" {...register('travel_style')}>
          <option value="relaxed">ゆったり</option>
          <option value="balanced">バランス</option>
          <option value="active">アクティブ</option>
        </select>
      </div>

      <button type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? 'プラン生成中...' : 'プランを生成'}
      </button>

      {mutation.isError && (
        <div className="error">エラー: {mutation.error.message}</div>
      )}
    </form>
  )
}
