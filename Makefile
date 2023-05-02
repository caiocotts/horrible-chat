run: assets/font/3270.ttf
	go run github.com/cespare/reflex@latest -s go run .

build: assets/font/3270.ttf
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

assets/font/3270.ttf:
	mkdir -p assets/font
	curl "https://cdn.discordapp.com/attachments/634976915974782979/1102714571246014504/IBM_3270_Nerd_Font_Complete.ttf" > assets/font/3270.ttf
