import { describe, it, expect } from 'vitest'
import { buildItineraryFilename } from './buildItineraryFilename'

describe('buildItineraryFilename', () => {
  it('returns filename with destination and date', () => {
    const result = buildItineraryFilename('Tokyo', '2025-07-01')

    expect(result).toBe('travelist-Tokyo-2025-07-01.json')
  })

  it('replaces spaces in destination with hyphens', () => {
    const result = buildItineraryFilename('New York', '2025-08-15')

    expect(result).toBe('travelist-New-York-2025-08-15.json')
  })

  it('removes characters unsafe for filenames', () => {
    const result = buildItineraryFilename('Tokyo/Osaka', '2025-09-01')

    expect(result).toBe('travelist-TokyoOsaka-2025-09-01.json')
  })

  it('handles Japanese destination names', () => {
    const result = buildItineraryFilename('京都', '2025-10-01')

    expect(result).toBe('travelist-京都-2025-10-01.json')
  })
})
