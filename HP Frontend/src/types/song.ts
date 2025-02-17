export interface Song {
  id: string
  uri: string
  name: string
  artists: string[]
  album: string
  image: Image
  duration_ms: number
  explicit: boolean
  externalUrl: string
}
interface Image {
  url: string
  height: number
  width: number
}
