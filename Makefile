db:
	docker-compose --file docker-compose.db.yml up -d

build:
	docker-compose --file docker-compose.build.yml build

push:
	docker-compose --file docker-compose.build.yml build
	docker login https://ghcr.io/
	docker push ghcr.io/konstfish/angler_ingress:latest
	docker push ghcr.io/konstfish/angler_geoip-api:latest
	docker push ghcr.io/konstfish/angler_frontend:latest
	docker push ghcr.io/konstfish/angler_auth:latest
	docker push ghcr.io/konstfish/angler_backend:latest

deploy:
	kubectl rollout restart deployment/ingress-deployment
	kubectl rollout restart deployment/frontend-deployment
	kubectl rollout restart deployment/geoip-api-deployment

redeploy: push deploy
