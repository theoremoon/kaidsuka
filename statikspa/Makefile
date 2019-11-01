.PHONY: build
build: build-ui go

.PNONY: build-ui
build-ui:
	cd ui;\
		yarn;\
		yarn run parcel build index.html;
	rm -rf ./dist
	mv ui/dist ./dist
	cp -r ui/assets ./dist/assets

.PHONY: go
go:
	go generate
	go build -o app


