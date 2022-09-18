# gek

GEnerate Kratos

## Project multirepo

- gitsite/company/project

  - back
    - back-service-1
    - back-service-2
    - back-service-3
    - .settings (file)
  - deploy
    - local
      - docker-compose.yaml (file)
      - services.yaml (file)
    - stage
    - test
    - prod
  - front
    - front-service-1
    - front-service-2
  - proto

back & front projects apply to common proto 

Workdir of gek is back with .settings file

Also gek works in local dir with services.yaml file
