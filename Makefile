include .env
MIGRATE= migrate
# @read -p  "What is the name of migration?" NAME; 

#Not Reading input rn because Read is not working
create:
	@echo "Testing"

	${MIGRATE} create -ext sql -seq -dir ./database/migration  test_question_table

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 
force:
	@read -p  "Which version do you want to force?" VERSION; \
	$(MIGRATE) force $$VERSION

goto:
	@read -p  "Which version do you want to migrate?" VERSION; \
	$(MIGRATE) goto $$VERSION

drop:
	$(MIGRATE) drop

.PHONY: migrate-up migrate-down force goto drop create

.PHONY: migrate-up migrate-down force goto drop create auto-create
