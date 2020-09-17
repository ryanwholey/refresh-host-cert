up:
	
	. ./.env && docker-compose up --build -d
	until nc -z localhost 8200 ; do sleep 1 ; done
	. ./.env && cd vault && terraform init && terraform apply -auto-approve

down:
	
	. ./.env && cd vault && terraform destroy -auto-approve
	. ./.env && docker-compose down
