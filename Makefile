up:
	. ./.env
	docker-compose up --build -d
	sleep 5
	cd vault && terraform apply -auto-approve

down:
	cd vault && terraform destroy -auto-approve
	docker-compose down
