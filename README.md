# UpSnap v3

## Stable version can be found on [v2 branch](https://github.com/seriousm4x/UpSnap/tree/v2)

> This is the dev branch for v3. It's very much in progress but will have a release candidate once all boxes below are checked.

## Help developing

Fork this branch and clone it.

1. Start backend

```sh
cd backend
go mod tidy
go run main.go serve
```

Log in to [http://127.0.0.1:8090/\_/](http://127.0.0.1:8090/_/), create an admin user and add some devices.

2. Start frontend

```sh
cd frontend
pnpm i
pnpm run dev
```

Open up [http://localhost:5173/](http://localhost:5173/)

## To do

- frontend

  - [x] ~~add rest of settings page~~
  - [x] ~~form for adding new devices~~
  - [x] ~~add per device settings~~
  - [x] ~~theme toggle button~~
  - [x] ~~add device ports to cards~~
  - [ ] settings: add network scan
  - [ ] settings: add option for website title

- backend

  - [x] ~~make sure ping works~~
  - [x] ~~make sure arp works~~
  - [x] ~~make sure wake works~~
  - [x] ~~add shutdown command~~
  - [x] ~~add scheduled wake~~
  - [x] [~~#34 Add support for WOL passwords~~](https://github.com/seriousm4x/UpSnap/issues/34)
  - [x] [~~#33 Seemingly High Ram usage~~](https://github.com/seriousm4x/UpSnap/issues/33)
  - [x] [~~#32 API available?~~](https://github.com/seriousm4x/UpSnap/issues/32)

- [x] ~~add db import from v2~~
- [ ] move docker images to github packages and setup new workflow file
- [ ] create dockerfile and docker-compose.yml
