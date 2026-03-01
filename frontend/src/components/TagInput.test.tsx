import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TagInput } from './TagInput'

const suggestions = ['文化', '歴史', '食事', '自然', 'アート']

describe('TagInput', () => {
  it('renders provided tags as chips', () => {
    render(
      <TagInput tags={['文化', '食事']} onChange={() => {}} suggestions={suggestions} />,
    )

    expect(screen.getByText('文化')).toBeInTheDocument()
    expect(screen.getByText('食事')).toBeInTheDocument()
  })

  it('adds tag when Enter is pressed', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={[]} onChange={onChange} suggestions={suggestions} />,
    )

    const input = screen.getByRole('textbox')
    await user.type(input, '温泉{Enter}')

    expect(onChange).toHaveBeenCalledWith(['温泉'])
  })

  it('does not add duplicate tag', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={['文化']} onChange={onChange} suggestions={suggestions} />,
    )

    const input = screen.getByRole('textbox')
    await user.type(input, '文化{Enter}')

    expect(onChange).not.toHaveBeenCalled()
  })

  it('removes tag when delete button is clicked', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={['文化', '食事']} onChange={onChange} suggestions={suggestions} />,
    )

    const removeButtons = screen.getAllByRole('button', { name: '削除' })
    await user.click(removeButtons[0])

    expect(onChange).toHaveBeenCalledWith(['食事'])
  })

  it('adds tag when suggestion chip is clicked', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={[]} onChange={onChange} suggestions={suggestions} />,
    )

    await user.click(screen.getByRole('button', { name: '歴史 を追加' }))

    expect(onChange).toHaveBeenCalledWith(['歴史'])
  })

  it('hides suggestions that are already added as tags', () => {
    render(
      <TagInput tags={['文化', '歴史']} onChange={() => {}} suggestions={suggestions} />,
    )

    expect(screen.queryByRole('button', { name: '文化 を追加' })).not.toBeInTheDocument()
    expect(screen.queryByRole('button', { name: '歴史 を追加' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: '食事 を追加' })).toBeInTheDocument()
  })

  it('trims whitespace from input before adding', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={[]} onChange={onChange} suggestions={suggestions} />,
    )

    const input = screen.getByRole('textbox')
    await user.type(input, '  温泉  {Enter}')

    expect(onChange).toHaveBeenCalledWith(['温泉'])
  })

  it('does not add empty tag', async () => {
    const onChange = vi.fn()
    const user = userEvent.setup()

    render(
      <TagInput tags={[]} onChange={onChange} suggestions={suggestions} />,
    )

    const input = screen.getByRole('textbox')
    await user.type(input, '{Enter}')

    expect(onChange).not.toHaveBeenCalled()
  })
})
