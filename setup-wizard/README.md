# Assistenten (setup-wizard)

Go-koden bag GF2 IT Setup-Assistenten: én kodebase, der kompileres til en
Windows-`.exe` og en Mac-binær. Binæren starter en lokal webserver på
localhost, åbner elevens browser og viser den danske tringuide.

## Tests

    go test ./...

Alle tests kører uden at røre det rigtige OS (OS-laget fakes via
`internal/osops/osfake`).

## Byg elev-binærerne

    GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o dist/Assistenten.exe ./cmd/assistent
    GOOS=darwin  GOARCH=arm64 go build -o dist/Assistenten ./cmd/assistent

`-H=windowsgui` forhindrer et sort konsolvindue bag browseren, når eleven
dobbeltklikker på .exe'en. Ingen CGO og ingen dependencies — bygget virker
fra enhver maskine med Go installeret.

## Kør lokalt under udvikling (WSL/Linux)

    go run ./cmd/assistent

Linux-implementeringen er en udviklerstub: Wi-Fi-status svarer "intet
netværk", og SketchUp-installation går altid til fallback. Hvis browseren
ikke åbner automatisk, står URL'en i terminalen.
