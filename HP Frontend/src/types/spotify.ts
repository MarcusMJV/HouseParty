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

export interface SpotifyPlayer {
  addListener(
    event: 'initialization_error' | 'authentication_error' | 'account_error' | 'playback_error' | 'player_state_changed' | 'ready' | 'not_ready',
    callback: (params: any) => void
  ): void;
  connect(): Promise<boolean>;
  disconnect(): void;
  // You can add more methods as needed.
}

export interface SpotifySDK {
  Player: new (options: {
    name: string;
    getOAuthToken: (cb: (token: string) => void) => void;
    volume?: number;
  }) => SpotifyPlayer;
}




