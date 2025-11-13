# tagyou

## presentation

yet another mqtt broker.

This is a personal project started to learn go language and empower my understanding of mqtt protocol in IOT environment.

Development is driven by :
* mqtt 5 specifications by OASIS
  https://mqtt.org/mqtt-specification/ 

* pure golang implementation

* http apis to configure and inspect behaviour

* ready for kubernetes deployment, prometheus friendly metrics exposed

* futuristic support for controlling security and control through machine learning

## build locally
there is a very simple Makefile in the project root.

* make build : install dependencies and build the project
* make clean : remove any user/build data
* make init : remove any user data to start with an empty db and new admin password

## first start
HTTP apis are always authenticated, unlike the mqtt clients. First user to be created is "admin".
To set your password for "admin", at first launch, pass INIT_ADMIN_PASSWORD as env var so a user with username "admin" is created with selected password and you can start use apis (to create more users? register clients ?). All users can access everything.

```
docker run -v tagyou_data:/db -e DB_PATH=/db -e INIT_ADMIN_PASSWORD=my_fantastic_secure_password ilgianlu/tagyou
```

Last command does not reset db, if you already have one, rerun same command to update admin password.
On next run db is set so you just need to attach your volume and have tagyou search db on it.

```
docker run -v tagyou_data:/db -e DB_PATH=/db ilgianlu/tagyou
```

Look at included docker compose file for other often used configuration vars.

## apis

### auth and calling

to authenticate and get an api token :

```
curl -v http://localhost:8080/auth -d '{"Username":"admin","InputPassword":"my_fantastic_secure_password"}'
```

response :

```
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiO....."}
```

Use Authorization header in api calls

```
curl -H "Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQ......" http://localhost:8080/clients
```

### available resources
* POST /auth
* GET /clients
* GET /clients/{id}
* POST /clients
* DELETE /clients/{id}
* GET /sessions
* GET /subscriptions
* GET /users
* POST /users
* DELETE /users/{id}

## deploy

The project is NOT READY for production deployment.
Can be tested quite easily pulling from docker hub

> docker pull ilgianlu/tagyou

Tag "latest" is continuously refresh by github actions built on main branch.

## contribution

I'm developing on linux, visual studio code, go lang 1.24.

Build and modify should not more difficult than cloning the repo, branching, opening a pull request. Please open an issue before contributing and feel free to include @ilgianlu
for the review.
