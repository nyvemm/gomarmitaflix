# Go marmitaflix

## Description

This microservice will search for an easy marmita then you will not starve.


## Routes

### GET: `/movies/search/marmita`

Search for movies with the word "marmita" in the title.


### GET: `/open/marmita_download_link`

Get the download link for the movie.

## Example Usage

### Request

```http
GET http://127.0.0.1:3000/movies/search/marvel
```

### Response

```json
[
  {
    "title": "Permalink to Baixar The Marvels Dublado BluRay 720p | 1080p (2023) Download",
    "slug": "baixar-the-marvels-dublado-77/",
    "image": "https://www.themoviedb.org/t/p/w600_and_h900_bestv2/gmtrDKIKXRFkChgHGALZTiKDVo0.jpg",
    "download_link": "http://127.0.0.1:3000/open/baixar-the-marvels-dublado-77/"
  },
  {
    "title": "Permalink to As Marvels Torrent (2023) Dublado Oficial / Legendado HDCAM 720p | 1080p Download",
    "slug": "as-marvels-torrent-2023/",
    "image": "https://image.tmdb.org/t/p/w342/sPmmgdmApfjX9x2mg02bo0aUOU9.jpg",
    "download_link": "http://127.0.0.1:3000/open/as-marvels-torrent-2023/"
  }
]

```

```mermaid
graph LR
    B -->|Marmita exists| C(Fabricio got happy)
    B -->|Marmita does not exist| D(Fabricio got hungry)
    A(BOLACHA DE TERRA) -->|Check in MARMITARIA| B[MARMITARIA]
```