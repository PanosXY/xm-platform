lint:
	golangci-lint run -v --exclude-use-default=false --timeout 1m0s

compose:
	docker-compose up

setupdb:
	psql -h localhost -U xmadmin -f setup_db.sql
