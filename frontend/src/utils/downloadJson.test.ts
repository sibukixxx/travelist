import { describe, it, expect, vi, afterEach } from 'vitest'
import { downloadJson } from './downloadJson'

describe('downloadJson', () => {
  const createObjectURLMock = vi.fn().mockReturnValue('blob:mock-url')
  const revokeObjectURLMock = vi.fn()

  let clickedAnchors: { href: string; download: string; click: ReturnType<typeof vi.fn> }[]

  afterEach(() => {
    vi.restoreAllMocks()
    clickedAnchors = []
  })

  function setup() {
    clickedAnchors = []
    global.URL.createObjectURL = createObjectURLMock
    global.URL.revokeObjectURL = revokeObjectURLMock

    vi.spyOn(document, 'createElement').mockImplementation((tag: string) => {
      if (tag === 'a') {
        const anchor = { href: '', download: '', click: vi.fn() }
        clickedAnchors.push(anchor)
        return anchor as unknown as HTMLAnchorElement
      }
      return document.createElement(tag)
    })
    vi.spyOn(document.body, 'appendChild').mockImplementation((node) => node)
    vi.spyOn(document.body, 'removeChild').mockImplementation((node) => node)
  }

  it('creates a Blob with JSON content and triggers download', () => {
    setup()
    const data = { destination: 'Tokyo', days: 3 }

    downloadJson(data, 'test-file.json')

    expect(createObjectURLMock).toHaveBeenCalledOnce()
    const blob = createObjectURLMock.mock.calls[0][0] as Blob
    expect(blob).toBeInstanceOf(Blob)
    expect(blob.type).toBe('application/json')
  })

  it('sets correct filename on anchor element', () => {
    setup()

    downloadJson({ key: 'value' }, 'my-plan.json')

    expect(clickedAnchors[0].download).toBe('my-plan.json')
  })

  it('clicks the anchor to trigger download', () => {
    setup()

    downloadJson({}, 'file.json')

    expect(clickedAnchors[0].click).toHaveBeenCalledOnce()
  })

  it('cleans up by revoking object URL and removing anchor', () => {
    setup()

    downloadJson({}, 'file.json')

    expect(revokeObjectURLMock).toHaveBeenCalledWith('blob:mock-url')
    expect(document.body.removeChild).toHaveBeenCalledOnce()
  })

  it('serializes data with 2-space indentation', () => {
    setup()
    const data = { a: 1 }

    downloadJson(data, 'file.json')

    const calls = createObjectURLMock.mock.calls
    const blob = calls[calls.length - 1][0] as Blob
    expect(blob.size).toBe(JSON.stringify(data, null, 2).length)
  })
})
