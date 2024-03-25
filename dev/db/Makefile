# DB Control

up:
	docker compose up -d

down:
	docker compose down

cleanup: down
	docker volume rm snippetbox_mysql_data

console:
	docker compose exec -it mysql mysql -usnippetbox -psnippetbox snippetbox


# Executing snippets

schemas: snippets-table session-table dummy-records

snippets-table:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox < /schemas/snippets-table.sql'

session-table:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox < /schemas/session-table.sql'

dummy-records:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox < /schemas/dummy-records.sql'

drop-table:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox -e "DROP TABLE snippets;"'

drop-records:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox -e " DELETE FROM snippets;"'


# Verifying

check-tables:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox -e "SHOW TABLES;"'

check-records:
	docker compose exec mysql sh -c 'mysql -usnippetbox -psnippetbox snippetbox -e "SELECT id, title, expires FROM snippets;"'


# DB Debug

logs:
	docker compose logs -f

shell:
	docker compose exec -it mysql bash

root:
	docker compose exec -it mysql mysql -uroot -psnippetbox
