export default async function copyToClipboard(s: string) {
  try {
    await navigator.clipboard.writeText(s)
    return true
  } catch (err) {
    return false
  }
}