up:
	docker-compose up -d
psql:
	docker exec -it postgres psql -U postgres -d force
test:
	test_status=0;\
	docker-compose -f docker-compose.test.yaml up --build -d;\
	docker-compose -f docker-compose.test.yaml run integration_test go test -v -count=1 ./...  || test_status=$$?;\
	docker-compose -f docker-compose.test.yaml down; echo "status="$$test_status;exit $$test_status;
