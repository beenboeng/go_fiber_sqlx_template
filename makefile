.PHONY : migrate up down


export DATABASE_URL=$(shell grep DB_URL .env | cut -d '=' -f2)

up: 
	goose -dir migrations postgres $(DATABASE_URL) up


down: 
	goose -dir migrations postgres $(DATABASE_URL) down
