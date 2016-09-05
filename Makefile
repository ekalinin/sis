NAME=sis
EXEC=${NAME}
GOVER=1.7
ENVNAME=${NAME}${GOVER}

build: deps
	go build -o ${EXEC} main.go

deps:
	go get github.com/julienschmidt/httprouter

run:
	./${EXEC} -pretty-json

run-test:
	wrk -t5 -c5 -d5s --latency http://localhost:8000/generate/png/100/100

run-test2:
	ab -n 5000 -c 10 http://localhost:8000/generate/png/100/100 &
	ab -n 5000 -c 10 http://localhost:8000/generate/jpg/200/200 &
	ab -n 5000 -c 10 http://localhost:8000/generate/jpg/300/300

show-stats:
	watch "curl -s http://localhost:8000/stats"


test:
	@go test -v

#
# For virtual environment create with
# https://github.com/ekalinin/envirius
#
env-create: env-init env-deps

env-init:
	@bash -c ". ~/.envirius/nv && nv mk ${ENVNAME} --go-prebuilt=${GOVER}"

env-build:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make build'"

env-deps:
	@bash -c ". ~/.envirius/nv && nv do ${ENVNAME} 'make deps'"

env:
	@bash -c ". ~/.envirius/nv && nv use ${ENVNAME}"

