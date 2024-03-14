# DB Control

up:
	docker compose up -d

down:
	docker compose down

console:
	docker compose exec -it mysql mysql -usnippetbox -psnippetbox

# DB Debug

logs:
	docker compose logs -f


cleanup: down
	docker volume rm snippetbox_mysql_data
