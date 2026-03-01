export function buildItineraryFilename(destination: string, startDate: string): string {
  const sanitized = destination
    .replace(/\s+/g, '-')
    .replace(/[/\\:*?"<>|]/g, '')
  return `travelist-${sanitized}-${startDate}.json`
}
