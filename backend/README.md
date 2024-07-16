# Music Library Management Backend

## APIs

### Music Library Management APIs

1. **Add a New Music Track with Cover Image and MP3 File**
   - **Endpoint:** `/api/tracks` (POST)
   - **Description:** Add a new music track with details like title, cover image, artist, album, genre, release year, duration, and upload the cover image and MP3 file in a single request.
   - **Request Body:**
     - Form data with the following fields:
       - `title` (string, required)
       - `cover_image` (file, required)
       - `artist` (string, required)
       - `album` (string, optional)
       - `genre` (string, optional)
       - `release_year` (integer, optional)
       - `duration` (integer, required)
       - `mp3_file` (file, required)
   - **Sample cURL Request:**
     ```bash
     curl -X POST http://localhost:8080/api/tracks -F "title=Bài hát mới" -F "cover_image=@/cover.jpg" -F "artist=Ca sĩ A" -F "album=Album B" -F "genre=Pop" -F "release_year=2021" -F "duration=240" -F "mp3_file=@/file.mp3"
     ```

2. **View Details of a Specific Music Track**
   - **Endpoint:** `/api/tracks/:id` (GET)
   - **Description:** View the details of a specific music track by its ID.
   - **Request Parameters:** `id` - The ID of the music track.
   - **Sample cURL Request:**
     ```bash
     curl -X GET http://localhost:8080/api/tracks/60c72b2f9b1d8b6e9f3e9f3e
     ```

3. **Update an Existing Music Track**
   - **Endpoint:** `/api/tracks/:id` (PUT)
   - **Description:** Update the details of an existing music track, including the cover image.
   - **Request Parameters:** `id` - The ID of the music track.
   - **Request Body:**
     - Form data with the following fields:
       - `title` (string, optional)
       - `cover_image` (file, optional)
       - `artist` (string, optional)
       - `album` (string, optional)
       - `genre` (string, optional)
       - `release_year` (integer, optional)
       - `duration` (integer, optional)
       - `mp3_file` (file, optional)
   - **Sample cURL Request:**
     ```bash
     curl -X PUT http://localhost:8080/api/tracks/60c72b2f9b1d8b6e9f3e9f3e -F "title=Bài hát cập nhật" -F "cover_image=@/updated_cover.jpg" -F "artist=Ca sĩ B" -F "album=Album C" -F "genre=Rock" -F "release_year=2022" -F "duration=300" -F "mp3_file=@/updated_file.mp3"
     ```

4. **Delete a Music Track**
   - **Endpoint:** `/api/tracks/:id` (DELETE)
   - **Description:** Delete a music track from the library.
   - **Request Parameters:** `id` - The ID of the music track.
   - **Sample cURL Request:**
     ```bash
     curl -X DELETE http://localhost:8080/api/tracks/60c72b2f9b1d8b6e9f3e9f3e
     ```

5. **List All Music Tracks**
   - **Endpoint:** `/api/tracks` (GET)
   - **Description:** Display a list of all music tracks in the library.
   - **Request Query Parameters:** 
     - `page` - The page number for pagination (default is 1).
     - `limit` - The number of items per page (default is 10).
   - **Sample cURL Request:**
     ```bash
     curl -X GET http://localhost:8080/api/tracks?page=1&limit=10
     ```

6. **Play/Pause an MP3 File of a Music Track**
   - **Endpoint:** `/api/tracks/:id/play` (POST)
   - **Description:** Play or pause the MP3 file of a specified music track.
   - **Request Parameters:** `id` - The ID of the music track.
   - **Request Body:**
     ```json
     {
       "action": "play"  // or "pause"
     }
     ```
   - **Sample cURL Request:**
     ```bash
     curl -X POST http://localhost:8080/api/tracks/60c72b2f9b1d8b6e9f3e9f3e/play -H "Content-Type: application/json" -d '{
       "action": "play"
     }'
     ```

7. **Create a New Playlist**
   - **Endpoint:** `/api/playlists` (POST)
   - **Description:** Create a new playlist with a given name.
   - **Request Body:**
     ```json
     {
       "name": "Danh sách phát mới"
     }
     ```
   - **Sample cURL Request:**
     ```bash
     curl -X POST http://localhost:8080/api/playlists -H "Content-Type: application/json" -d '{
       "name": "Danh sách phát mới"
     }'
     ```

8. **Add a Track to a Playlist**
    - **Endpoint:** `/api/playlists/:playlistId/tracks/:trackId` (POST)
    - **Description:** Add a music track to a specified playlist.
    - **Request Parameters:**
      - `playlistId`: The ID of the playlist.
      - `trackId`: The ID of the music track.
    - **Sample cURL Request:**
      ```bash
      curl -X POST http://localhost:8080/api/playlists/60c72b2f9b1d8b6e9f3e9f3e/tracks/60c72b2f9b1d8b6e9f3e9f3e
      ```

9. **List All Playlists**
    - **Endpoint:** `/api/playlists` (GET)
    - **Description:** Display a list of all playlists.
    - **Request Query Parameters:** 
      - `page` - The page number for pagination (default is 1).
      - `limit` - The number of items per page (default is 10).
    - **Sample cURL Request:**
      ```bash
      curl -X GET http://localhost:8080/api/playlists?page=1&limit=10
      ```

10. **Search for Music Tracks and Playlists**
    - **Endpoint:** `/api/search` (GET)
    - **Description:** Search for music tracks and playlists by title, artist, album, or genre.
    - **Request Query Parameters:** 
      - `query` - The search query string.
      - `page` - The page number for pagination (default is 1).
      - `limit` - The number of items per page (default is 10).
    - **Sample cURL Request:**
      ```bash
      curl -X GET http://localhost:8080/api/search?query=B%C3%A0i%20h%C3%A1t&page=1&limit=10
      ```

11. **List All Genres**
    - **Endpoint:** `/api/genres` (GET)
    - **Description:** Provides a list of available genres.
    - **Sample cURL Request:**
      ```bash
      curl -X GET http://localhost:8080/api/genres
      ```
