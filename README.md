# House Party 🎉🎶

House Party is a collaborative music-listening experience where users can connect their Spotify account, create or join a room, and queue songs to a shared playlist in real-time. Think of it as a live, crowd-controlled radio station where everyone in the room helps shape the vibe.

> 🚀 **Live Demo:** [hp-frontend.up.railway.app/signup-or-login](https://hp-frontend.up.railway.app/)  
> 👤 **Demo Account:**  
> &nbsp;&nbsp;&nbsp;&nbsp;**Username:** `demoUser`  
> &nbsp;&nbsp;&nbsp;&nbsp;**Password:** `1234`
>
> 🛈 **Please note:** Since this app currently uses a **free Spotify developer account**, only pre-approved users added via the Spotify Developer Dashboard can play music through the app.
>
> - The `demoUser` account is the only account added for playback and acts as the **host device**.
> - If you log in with your own Spotify account, you can explore the app and **queue songs**, but **you won't hear any music** unless you're logged in as `demoUser`.
> - This limitation is in place for **proof-of-concept purposes** and will be lifted in future versions with proper production setup.

## 🔑 Features

- ✅ Spotify OAuth login & playback control
- 🏠 Create and manage private or public rooms
- 🚪 Join rooms via code or link
- 🎶 Add songs to a shared playlist queue
- ⏭️ Vote-based skipping (majority rule)
- 🔄 Real-time sync across users in a room

## 🛠️ Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** Vue.js
- **Spotify API:** For authentication and playback
- **WebSockets:** For real-time communication
- **Database:** SQLite (development only)

## 🚧 Disclaimer

> ⚠️ **Note:** House Party is currently a **proof of concept** and is actively in development.
>
> - The **live demo** is for **testing purposes only**.
> - The interface is **not yet optimized for mobile**.
> - I use **SQLite** in development (subject to change in production).
> - At the moment, **only the host device (the room creator)** plays the music. Other users contribute to the shared playlist but do **not play audio on their own devices**.
> - This setup is intended for **live social events** with a central speaker setup.
> - In future releases, the plan is to support **fully online synchronized playback across multiple devices**, allowing all participants to listen simultaneously in real time.
