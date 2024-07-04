override port = 3000

.PHONY: clean-docker-images

run:
	@echo "Run Go Application"
	@kill -9 $$(lsof -t -i:${port}) || true
	go run cmd/main.go

tidy:
	@echo "Install Packages"
	go mod tidy

rebuild-app:
	docker image prune -f && docker compose up --build myserver

down:
	docker compose down && docker image prune -f
		
ps:
	docker ps -a
	docker compose ps -a

up:
	docker compose up

rebuild-all:
	docker compose down && docker image prune -f && docker compose build --no-cache && docker compose up

clean:
	@echo "Removing all Docker images"
	@docker rmi $$(docker images -q) || true