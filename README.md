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

## first start
http apis are always authenticated, unlike the mqtt clients. First user to be created is "admin".
To set your password for "admin", at first launch, pass INIT_ADMIN_PASSWORD as env var so a user with username "admin" is created with selected password and you can start use apis (to create more users? register clients ?). All users can access everything.

```
docker run -e INIT_ADMIN_PASSWORD=my_fantastic_secure_password ilgianlu/tagyou
```

## deploy

The project is NOT READY for production deployment.
Can be tested quite easily pulling from docker hub

> docker pull ilgianlu/tagyou

Tag "latest" is continuously refresh by github actions built on main branch.

## contribution

I'm developing on linux, visual studio code, go lang 1.20.

Build and modify should not more difficult than cloning the repo, branching, opening a pull request. Please open an issue before contributing and feel free to include @ilgianlu
for the review.
