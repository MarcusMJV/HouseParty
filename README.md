# House Party 🎉🎶

House Party is a collaborative music-listening experience where users can connect their Spotify account, create or join a room, and queue songs to a shared playlist in real-time. Think of it as a live, crowd-controlled radio station where everyone in the room helps shape the vibe.

> 🚀 **Live Demo:** [hp-frontend.up.railway.app/signup-or-login](https://hp-frontend.up.railway.app/signup-or-login)  
> 👤 **Demo Account:**  
> &nbsp;&nbsp;&nbsp;&nbsp;**Username:** `demoUser`  
> &nbsp;&nbsp;&nbsp;&nbsp;**Password:** `1234`
>
> 🛈 **Please note:** You must log in with the **demo account** to hear music playback in the demo room, because that account is currently the **host device**.
>
> - You’re welcome to create your own account and connect your Spotify account.
> - If you create your own room, your device will become the host and playback will work on your account.
> - If you **join the demo room with a different account**, you will be able to **queue songs**, but you **won’t hear any music** as playback only happens on the host device.

## 🔑 Features

- ✅ Spotify OAuth login & playback control
- 🏠 Create and manage private or public rooms
- 🚪 Join rooms via code or link
- 🎶 Add songs to a shared playlist queue
- ⏭️ Vote-based skipping (majority rule)
- 🔄 Real-time sync across users in a room

## 🛠️ Tech Stack

- **Backend:** Go (Golang)
- **Frontend:** [your frontend tech, e.g., Vue.js / React]
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
