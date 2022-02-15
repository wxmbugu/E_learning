# e-learning
In progress

### Few Notes Swagger
  - #### Setup swagger
    - ##### Install swagger
        -  docker pull quay.io/goswagger/swagger
        -  alias swagger='sudo docker run --rm -it  --user $(id -u):$(id -g) -p 5000:5000 -e GOPATH=$(go env GOPATH):/go -e GOCACHE=/tmp -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger'
    - check version of swagger -> swagger version
    - generate swagger.json file -> swagger generate spec -o ./swagger.json
    - swagger serve ->  swagger serve ./swagger.json --no-open --port 5000
    - swagger serve documentation ->  swagger serve -F swagger ./swagger.json --no-open --port 5000



#### Todo
 -  Mock db
 -  Test Api
 - Setup Video Upload endpoint