.Phony: build up ps rm st

build:
	docker compose build api bot --no-cache

up:
	docker compose up -d

ps:
	docker compose ps

rm:
	docker compose down -v

ex:
	docker stop $$(docker ps -q)

st: build up ps

rt: rm st