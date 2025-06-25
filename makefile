include .env

run:
	$(MAKE) clean
	$(MAKE) build
	docker-compose up --build magasin logistique mere dbmagasin dblogistique dbMere

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/magasin/ ./caisse_app_scaled/magasin/app.go
	cp -rf caisse_app_scaled/magasin/view out/magasin/
	cp -rf caisse_app_scaled/commonjs out/magasin/commonjs/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/centre_logistique/ ./caisse_app_scaled/centre_logistique/app.go
	cp -rf caisse_app_scaled/centre_logistique/view out/centre_logistique/
	cp -rf caisse_app_scaled/commonjs out/centre_logistique/commonjs/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/maison_mere/ ./caisse_app_scaled/maison_mere/app.go
	cp -rf caisse_app_scaled/maison_mere/view out/maison_mere/
	cp -rf caisse_app_scaled/commonjs out/maison_mere/commonjs/

buildMagasin:
	rm -rf out/magasin/
	go build -o out/magasin/ ./caisse_app_scaled/magasin/app.go
	cp -rf caisse_app_scaled/magasin/view out/magasin/
buildLogistique:
	rm -rf out/centre_logistique/
	go build -o out/centre_logistique/ ./caisse_app_scaled/centre_logistique/app.go
	cp -rf caisse_app_scaled/centre_logistique/view out/centre_logistique/
buildMere:
	rm -rf out/maison_mere/
	go build -o out/maison_mere/ ./caisse_app_scaled/maison_mere/app.go
	cp -rf caisse_app_scaled/maison_mere/view out/maison_mere/

runMagasin:
	$(MAKE) buildMagasin
	cd out/magasin && DB_PORT=5432 GATEWAY=localhost DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app
runLogistique:
	$(MAKE) buildLogistique
	cd out/centre_logistique && DB_PORT=5433 GATEWAY=localhost DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app
runMere:
	$(MAKE) buildMere
	cd out/maison_mere && DB_PORT=5434 GATEWAY=localhost DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app

run-local:
	$(MAKE) runMere &
	$(MAKE) runLogistique &
	$(MAKE) runMagasin
	killall runMere
	killall runLogistique

test:
	$(MAKE) clean
	$(MAKE) comment-swagger
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run
	docker-compose up -d --build dbtest
	go test -v ./tests/...
	$(MAKE) unComment-swagger

docs: .FORCE
	go get github.com/swaggo/swag/cmd/swag@latest
	cd caisse_app_scaled && go run github.com/swaggo/swag/cmd/swag init -d magasin,models,utils -o ../docs/swagger/magasin -g api/api_data.go
	cd caisse_app_scaled && go run github.com/swaggo/swag/cmd/swag init -d centre_logistique,models,utils -o ../docs/swagger/logistique -g api/api_data.go
	cd caisse_app_scaled && go run github.com/swaggo/swag/cmd/swag init -d maison_mere,models,utils -o ../docs/swagger/mere -g api/api_data.go

clean:
	rm -rf out/*
	docker stop $$(docker ps -a -q) || true
	docker rm $$(docker ps -a -q) || true
# docker image rm $$(docker image ls -q) || true
# docker volume rm $$(docker volume ls -q) || true

dev-setup:
	echo $(PWD) | docker login -u $(USERNAME) --password-stdin
	go mod tidy
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(MAKE) docs

comment-swagger:
# ubuntu
	find caisse_app_scaled -type f -exec sed -i 's|^\(\s*_ "caisse-app-scaled/docs/swagger/logistique"\)|// \1|' {} + ||true
# mac
	find caisse_app_scaled -type f -exec sed -i '' '/_ "caisse-app-scaled\/docs\/swagger/s/^/\/\//g' {} + ||true

unComment-swagger:
# ubuntu
	find caisse_app_scaled -type f -exec sed -i 's|^// \(\s*_ "caisse-app-scaled/docs/swagger/logistique"\)|\1|' {} + ||true
# mac
	find caisse_app_scaled -type f -exec sed -i '' '/_ "caisse-app-scaled\/docs\/swagger/s/^\/\///g' {} + ||true

.FORCE:
