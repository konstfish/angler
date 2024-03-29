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
	kubectl -n angler rollout restart deployment/ingress-deployment
	kubectl -n angler rollout restart deployment/frontend-deployment
	kubectl -n angler rollout restart deployment/geoip-api-deployment
	kubectl -n angler rollout restart deployment/auth-deployment
	kubectl -n angler rollout restart deployment/backend-deployment

redeploy: push deploy
