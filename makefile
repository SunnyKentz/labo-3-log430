include .env

run:
	$(MAKE) clean-docker
	$(MAKE) build
	docker-compose up --build
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
build:
	rm -rf out/magasin/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/magasin/ ./caisse_app_scaled/magasin/app.go
	cp -rf caisse_app_scaled/magasin/view out/magasin/

	rm -rf out/centre_logistique/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/centre_logistique/ ./caisse_app_scaled/centre_logistique/app.go
	cp -rf caisse_app_scaled/centre_logistique/view out/centre_logistique/

	rm -rf out/maison_mere/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o out/maison_mere/ ./caisse_app_scaled/maison_mere/app.go
	cp -rf caisse_app_scaled/maison_mere/view out/maison_mere/
runMagasin:
	$(MAKE) buildMagasin
	cd out/magasin && DB_PORT=5432 DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app
runLogistique:
	$(MAKE) buildLogistique
	cd out/centre_logistique && DB_PORT=5433 DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app
runMere:
	$(MAKE) buildMere
	cd out/maison_mere && DB_PORT=5434 DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) ./app

test:
	$(MAKE) clean-docker
	golangci-lint run
	docker-compose up -d --build dbtest
	go test -v ./tests/...

clean:
	rm -rf out/*
	
clean-docker:
	docker stop $$(docker ps -a -q) || true
	docker rm $$(docker ps -a -q) || true
# docker image rm $$(docker image ls -q) || true
# docker volume rm $$(docker volume ls -q) || true

deploy:
	$(MAKE) runMagasin
	$(MAKE) runLogistique
	$(MAKE) runMere
	docker build -t $(USERNAME)/caisse-app-Magasin:latest ./caisse_app_scaled/magasin
	docker build -t $(USERNAME)/caisse-app-Logistique:latest ./caisse_app_scaled/centre_logistique
	docker build -t $(USERNAME)/caisse-app-Mere:latest ./caisse_app_scaled/maison_mere
# docker push $(USERNAME)/caisse-app-Magasin:latest
# docker push $(USERNAME)/caisse-app-Logistique:latest
# docker push $(USERNAME)/caisse-app-Mere:latest


dev-setup:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	echo $(PWD) | docker login -u $(USERNAME) --password-stdin
	go mod tidy

tag:
	@read -p "Enter tag to push: " tag; git push origin $$tag
	