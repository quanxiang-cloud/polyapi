name=polyapi_job
version=v0.7.3
devHost=192.168.200.20
devUser=ubuntu
repository=lowcode
dockerHost=dockerhub.qingcloud.com

env:
#-- open go mod vendor --
	go mod vendor

docker-test: env
	cd ../../.././ && \
	docker build -f ./pkg/jobs/jobmain/Dockerfile -t  $(dockerHost)/$(repository)/$(name):$(version) .
	docker push  $(dockerHost)/$(repository)/$(name):$(version)
