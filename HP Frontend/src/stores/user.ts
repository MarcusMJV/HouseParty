import { defineStore } from 'pinia';

interface UserCredentials {
  id: number;
  username: string;
  email: string;
  spotify_connected: boolean;
}

interface UserState {
  jwt: string | null;
  credentials: UserCredentials ;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    jwt: localStorage.getItem('jwt') || null,
    credentials: JSON.parse(localStorage.getItem('credentials') || 'null')
  }),
  getters: {
    isAuthenticated(): boolean {
      return !!this.jwt;
    }
  },
  actions: {
    setJwt(jwt: string) {
      this.jwt = jwt;
      localStorage.setItem('jwt', jwt);
    },
    activateSpotifyConnection() {
      this.credentials.spotify_connected = true
      localStorage.setItem('credentials', JSON.stringify(this.credentials));
    },
    setCredentials(credentials: UserCredentials) {
      this.credentials = credentials;
      localStorage.setItem('credentials', JSON.stringify(credentials));
    },
    clearUser() {
      this.jwt = null;
      localStorage.removeItem('jwt');
      localStorage.removeItem('credentials');
    },
  }
});
