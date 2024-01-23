pg_mig:
	migrate create -ext sql -dir postgres/migration/ -seq ${name}

