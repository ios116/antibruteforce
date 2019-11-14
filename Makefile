up:
	docker-compose up -d
psql:
	docker exec -it postgres psql -U postgres -d force
