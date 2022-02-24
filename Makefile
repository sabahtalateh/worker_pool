# goose -dir migrations postgres "host=127.0.0.1 port=23345 user=postgres password=postgres dbname=worker_pool sslmode=disable binary_parameters=yes" up
migrations-apply:
	@goose -dir migrations/pg postgres "${DB_URI}" up