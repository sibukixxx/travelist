import { useState, type KeyboardEvent } from 'react'

interface TagInputProps {
  tags: string[]
  onChange: (tags: string[]) => void
  suggestions: string[]
  placeholder?: string
}

export function TagInput({ tags, onChange, suggestions, placeholder = '入力してEnterで追加' }: TagInputProps) {
  const [input, setInput] = useState('')

  const addTag = (tag: string) => {
    const trimmed = tag.trim()
    if (trimmed === '' || tags.includes(trimmed)) return
    onChange([...tags, trimmed])
    setInput('')
  }

  const removeTag = (index: number) => {
    onChange(tags.filter((_, i) => i !== index))
  }

  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      addTag(input)
    } else if (e.key === 'Backspace' && input === '' && tags.length > 0) {
      removeTag(tags.length - 1)
    }
  }

  const handleInput = (value: string) => {
    if (value.includes(',')) {
      const parts = value.split(',')
      const tagToAdd = parts[0]
      addTag(tagToAdd)
      setInput(parts.slice(1).join(','))
    } else {
      setInput(value)
    }
  }

  const unusedSuggestions = suggestions.filter((s) => !tags.includes(s))

  return (
    <div className="tag-input-container">
      <div className="tag-list">
        {tags.map((tag, i) => (
          <span key={tag} className="tag">
            {tag}
            <button
              type="button"
              className="tag-remove"
              aria-label="削除"
              onClick={() => removeTag(i)}
            >
              ×
            </button>
          </span>
        ))}
        <input
          type="text"
          className="tag-input"
          value={input}
          onChange={(e) => handleInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder={tags.length === 0 ? placeholder : ''}
        />
      </div>
      {unusedSuggestions.length > 0 && (
        <div className="tag-suggestions">
          {unusedSuggestions.map((s) => (
            <button
              key={s}
              type="button"
              className="tag-suggestion"
              aria-label={`${s} を追加`}
              onClick={() => addTag(s)}
            >
              {s}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
