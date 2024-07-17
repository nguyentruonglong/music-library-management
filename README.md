# Music Library Management Website


### Project Directory Structure

```
music-library-management/
├── frontend/
│   ├── public/
│   │   ├── index.html
│   │   └── ...
│   ├── src/
│   │   ├── assets/
│   │   │   └── images/
│   │   ├── components/
│   │   │   ├── MusicTrack/
│   │   │   │   ├── AddTrackForm.jsx
│   │   │   │   ├── EditTrackForm.jsx
│   │   │   │   ├── TrackDetails.jsx
│   │   │   │   ├── TrackList.jsx
│   │   │   │   └── ...
│   │   │   ├── Playlist/
│   │   │   │   ├── AddPlaylistForm.jsx
│   │   │   │   ├── EditPlaylistForm.jsx
│   │   │   │   ├── PlaylistDetails.jsx
│   │   │   │   ├── PlaylistList.jsx
│   │   │   │   └── ...
│   │   │   ├── Player/
│   │   │   │   ├── PlayerControls.jsx
│   │   │   │   └── ...
│   │   │   ├── Search/
│   │   │   │   ├── SearchBar.jsx
│   │   │   │   └── SearchResults.jsx
│   │   │   └── ...
│   │   ├── context/
│   │   │   ├── MusicContext.jsx
│   │   │   └── PlaylistContext.jsx
│   │   ├── hooks/
│   │   │   ├── useMusic.jsx
│   │   │   └── usePlaylist.jsx
│   │   ├── pages/
│   │   │   ├── Home.jsx
│   │   │   ├── MusicTrackPage.jsx
│   │   │   ├── PlaylistPage.jsx
│   │   │   ├── SearchPage.jsx
│   │   │   └── ...
│   │   ├── services/
│   │   │   ├── api.js
│   │   │   └── ...
│   │   ├── App.jsx
│   │   ├── index.jsx
│   │   └── ...
│   ├── package.json
│   ├── .env.development
│   ├── .env.production
│   ├── README.md
│   └── ...
├── backend/
│   ├── api/
│   │   ├── controllers/
│   │   │   ├── track_controller.go
│   │   │   ├── playlist_controller.go
│   │   │   ├── genre_controller.go.go
│   │   │   └── ...
│   │   ├── models/
│   │   │   ├── track.go
│   │   │   ├── playlist.go
│   │   │   ├── genre.go
│   │   │   └── ...
│   │   ├── routes/
│   │   │   ├── track_routes.go
│   │   │   ├── playlist_routes.go
│   │   │   ├── genre_routes.go
│   │   │   └── ...
│   │   ├── services/
│   │   │   ├── track_service.go
│   │   │   ├── playlist_service.go
│   │   │   ├── genre_service.go
│   │   │   └── ...
│   │   ├── utils/
│   │   │   ├── db.go
│   │   │   ├── response.go
│   │   └── ...
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── ...
│   ├── uploads/
│   ├── config/
│   │   ├── config.go
│   ├── errors/
│   │   ├── handler.go
│   │   ├── messages.go
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   ├── .env.development
│   ├── .env.production
│   ├── README.md
│   └── ...
├── docker-compose.yml
├── Dockerfile
├── .gitignore
├── .dockerignore
└── README.md
```