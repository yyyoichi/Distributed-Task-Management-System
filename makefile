restart-%:
	docker-compose restart ${@:restart-%=%}

start:
	docker-compose start

stop:
	docker-compose stop