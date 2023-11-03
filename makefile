up:
	docker-compose up -d --build
exec:
	@echo "Run 'go run .' in the docker container"
	docker exec -w /workspace -it cli-app /bin/bash

logs-%:
	docker logs -f ${@:logs-%=%}

up-%:
	docker-compse ${@:up-%=%}

restart:
	docker-compose restart
restart-%:
	docker-compose restart ${@:restart-%=%}

start:
	docker-compose start
start-%:
	docker-compse start ${@:start-%=%}

stop:
	docker-compose stop
stop-%:
	docker-compose stop ${@:stop-%=%}

rm:
	docker-compose rm 
rm-%:
	docker-compose rm ${@:rm-%=%}