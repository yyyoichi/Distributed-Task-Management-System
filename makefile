restart-%:
	docker-compose restart ${@:restart-%=%}

logs-%:
	docker logs -f ${@:logs-%=%}

start:
	docker-compose start

stop:
	docker-compose stop