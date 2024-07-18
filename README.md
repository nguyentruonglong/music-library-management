# Music Library Management Website

## Description

The Music Library Management system is a full-stack web application that allows users to organize and manage their music collections, including tracks and playlists. Users can add, update, search, and delete music tracks, upload mp3 files, view details of specific tracks, and manage playlists. The application provides a user-friendly interface with ReactJS and Ant Design for the frontend, and uses Golang with the Gin framework for the backend. MongoDB is used as the database.


### Project Directory Structure

```
music-library-management/
├── frontend/                               # Frontend source code
│   ├── public/                             # Static public assets
│   │   ├── index.html                      # Main HTML file
│   │   └── ...                  
│   ├── src/                                # React application source code
│   │   ├── assets/                         # Asset files such as images
│   │   │   └── images/                     # Image files
│   │   ├── components/                     # Reusable React components
│   │   │   ├── MusicTrack/                 # Components related to music tracks
│   │   │   │   ├── AddTrackForm.jsx        # Form to add a new track
│   │   │   │   ├── EditTrackForm.jsx       # Form to edit an existing track
│   │   │   │   ├── TrackDetails.jsx        # Display details of a track
│   │   │   │   ├── TrackList.jsx           # List of music tracks
│   │   │   │   └── ...          
│   │   │   ├── Playlist/                   # Components related to playlists
│   │   │   │   ├── AddPlaylistForm.jsx     # Form to add a new playlist
│   │   │   │   ├── EditPlaylistForm.jsx    # Form to edit an existing playlist
│   │   │   │   ├── PlaylistDetails.jsx     # Display details of a playlist
│   │   │   │   ├── PlaylistList.jsx        # List of playlists
│   │   │   │   └── ...          
│   │   │   ├── Player/                     # Components for music player controls
│   │   │   │   ├── PlayerControls.jsx      # Controls for the music player
│   │   │   │   └── ...          
│   │   │   ├── Search/                     # Components for search functionality
│   │   │   │   ├── SearchBar.jsx           # Search bar component
│   │   │   │   └── SearchResults.jsx       # Display search results
│   │   │   └── ...              
│   │   ├── context/                        # Context providers for state management
│   │   │   ├── MusicContext.jsx            # Context for music tracks
│   │   │   └── PlaylistContext.jsx         # Context for playlists
│   │   ├── hooks/                          # Custom hooks for reusable logic
│   │   │   ├── useMusic.jsx                # Hook for music-related logic
│   │   │   └── usePlaylist.jsx             # Hook for playlist-related logic
│   │   ├── pages/                          # Page components for different routes
│   │   │   ├── Home.jsx                    # Home page
│   │   │   ├── MusicTrackPage.jsx          # Page for managing music tracks
│   │   │   ├── PlaylistPage.jsx            # Page for managing playlists
│   │   │   ├── SearchPage.jsx              # Page for search results
│   │   │   └── ...              
│   │   ├── services/                       # Services for API calls and business logic
│   │   │   ├── api.jsx                     # API service for HTTP requests
│   │   │   └── ...              
│   │   ├── App.jsx                         # Main App component
│   │   ├── index.jsx                       # Entry point for the React application
│   │   └── ...                  
│   ├── package.json                        # Project dependencies and scripts
│   ├── .env.development                    # Environment variables for development
│   ├── .env.production                     # Environment variables for production
│   ├── README.md                           # Frontend documentation
│   └── ...                      
├── backend/                                # Backend source code
│   ├── api/                                # API-related code
│   │   ├── controllers/                    # API endpoint controllers
│   │   │   ├── track_controller.go         # Controller for track endpoints
│   │   │   ├── playlist_controller.go      # Controller for playlist endpoints
│   │   │   ├── genre_controller.go         # Controller for genre endpoints
│   │   │   ├── search_controller.go        # Controller for search endpoints
│   │   │   ├── file_controller.go          # Controller for file endpoints
│   │   │   └── ...              
│   │   ├── models/                         # Database models
│   │   │   ├── track.go                    # Model for track
│   │   │   ├── playlist.go                 # Model for playlist
│   │   │   ├── genre.go                    # Model for genre
│   │   │   ├── file.go                     # Model for file metadata
│   │   │   └── ...              
│   │   ├── routes/                         # API routes
│   │   │   ├── track_routes.go             # Routes for track endpoints
│   │   │   ├── playlist_routes.go          # Routes for playlist endpoints
│   │   │   ├── genre_routes.go             # Routes for genre endpoints
│   │   │   ├── search_routes.go            # Routes for search endpoints
│   │   │   ├── file_routes.go              # Routes for file endpoints
│   │   │   └── ...              
│   │   ├── services/                       # Business logic and services
│   │   │   ├── track_service.go            # Service for track logic
│   │   │   ├── playlist_service.go         # Service for playlist logic
│   │   │   ├── genre_service.go            # Service for genre logic
│   │   │   ├── search_service.go           # Service for search logic
│   │   │   ├── file_service.go             # Service for file management
│   │   │   └── ...              
│   │   ├── utils/                          # Utility functions and helpers
│   │   │   ├── db.go                       # Database connection and setup
│   │   │   ├── response.go                 # Helper for standardized responses
│   │   │   ├── seed.go                     # Database seeding script
│   │   │   ├── url.go                      # Helpers that assist in working with URLs
│   │   └── ...              
│   ├── middleware/                         # Middleware for request processing
│   │   ├── auth_middleware.go              # Authentication middleware
│   │   └── ...              
│   ├── uploads/                            # Directory for uploaded files
│   ├── config/                             # Configuration files
│   │   ├── config.go                       # Application configuration
│   ├── errors/                             # Error handling
│   │   ├── handler.go                      # Centralized error handler
│   │   ├── messages.go                     # Error messages
│   ├── main.go                             # Entry point for the backend server
│   ├── go.mod                              # Go module dependencies
│   ├── go.sum                              # Go module checksums
│   ├── .env.development                    # Environment variables for development
│   ├── .env.production                     # Environment variables for production
│   ├── README.md                           # Backend documentation
│   └── ...                             
├── docker-compose.yml                      # Docker Compose configuration
├── Dockerfile                              # Dockerfile for building the image
├── .gitignore                              # Git ignore file
├── .dockerignore                           # Docker ignore file
└── README.md                               # Project documentation
```

## Prerequisites

Ensure you have the following software installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Node.js](https://nodejs.org/)

## Environment Variables

Set up your environment variables in the following files:

- `.env.development` for development environment variables
- `.env.production` for production environment variables

## Running the Application with Docker

### Build the Docker images

To build the Docker images for the frontend and backend services, run the following command:

```bash
docker-compose -f docker-compose.yml build
```

### Start the Docker containers

To start the application in detached mode, run the following command:

```bash
$ docker-compose -f docker-compose.yml up -d
```

This command will start the backend and MongoDB services.

### Stopping the Docker containers

To stop the running containers, run the following command:

```bash
$ docker-compose -f docker-compose.yml down
```

### Accessing the Application

- Frontend: Open your web browser and navigate to http://localhost:3000
- Backend: The backend API will be accessible at http://localhost:8080/api

### Accessing MongoDB

To access MongoDB using MongoDB Compass or any other MongoDB client, connect to:

```bash
$ mongodb://localhost:27017/musiclibrary
```

## Contributing

Contributions are welcome! Please create a new issue or submit a pull request to contribute.

## Credits

- Author: Nguyen Truong Long

## Contact

If you have any questions or need further assistance, feel free to contact me.

This [README.md](README.md) file provides a detailed guide to setting up and running your Music Library Management system using Docker, along with other conventional sections. Adjust the placeholders me and my email with your actual contact information.
