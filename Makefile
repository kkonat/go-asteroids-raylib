loader: cmd\loader\main.go
	cd cmd\loader && go build
	mv cmd\loader\loader.exe loader\loader.exe

game: main.go loader
	set CGO_ENABLED=1&& go build -ldflags "-s -w -linkmode 'external' -extldflags '-static'"
game.zip: game game.zip
	7z a -tzip -mx=9 game.zip bangbang.exe
	python renamezip.py
clean:
	rm .loader\loader.exe
