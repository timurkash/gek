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

Workdir for gek is `back` with `.settings` file

This utility can be used as tool for business-oriented resource development

To install run
```
go install github.com/timurkash/gek@latest
```

The development methodology is described in https://www.youtube.com/watch?v=5OkPz78Xk3M (russian)
