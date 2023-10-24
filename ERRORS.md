# NVI API Error Codes

Click here for [field validation rules](https://github.com/go-playground/validator#baked-in-validations)

## E000

- Error Name: `Unknown`
- Status: `500`
- Explanation: the server ran into an unexpected error, see server logs or response error for more details

## E001

- Error Name: `RegisterInvalidBody`
- Controller: `user`
- Path: `/register`
- Method: `POST`
- Status: `400`
- Content: `application/json`
- Body: `name, email, password`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - name: `required,gte=2,lte=64`
    - email: `required,email,lte=255`
    - password: `required,gte=5,lte=36`

## E002

- Error Name: `RegisterEmailTaken`
- Controller: `user`
- Path: `/register`
- Method: `POST`
- Status: `200`
- Content: `application/json`
- Body: `name, email, password`
- Explanation: the request body contains an `email` field that's already in use


## E003

- Error Name: `LoginInvalidBody`
- Controller: `user`
- Path: `/login`
- Method: `POST`
- Status: `400`
- Content: `application/json`
- Body: `email, password`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - email: `required,email,lte=255`
    - password: `required,gte=5,lte=36`


## E004

- Error Name: `LoginUnregisteredEmail`
- Controller: `user`
- Path: `/login`
- Method: `POST`
- Status: `200`
- Content: `application/json`
- Body: `email, password`
- Explanation: the request body contains an unregistered `email` field

## E005

- Error Name: `LoginInvalidPassword`
- Controller: `user`
- Path: `/login`
- Method: `POST`
- Status: `401`
- Content: `application/json`
- Body: `email, password`
- Explanation: the request body contains an invalid `password` field for the provided `email` field

## E006

- Error Name: `LoginAccountNotVerified`
- Controller: `user`
- Path: `/login`
- Method: `POST`
- Status: `401`
- Content: `application/json`
- Body: `email, password`
- Explanation: the request body contains an `email` field that hasn't been verified yet

## E007

- Error Name: `VerifyAccountInvalidToken`
- Controller: `user`
- Path: `/verify/account`
- Method: `PATCH`
- Status: `401`
- Query: `token`
- Explanation: a `token` that was assigned as a query `?token=` is invalid (missing, expired or wrong signature); another
token may need to be regenerated

## E008

- Error Name: `ResendAccountVerificationInvalidEmail`
- Controller: `user`
- Path: `/reverify/account`
- Method: `PATCH`
- Status: `400`
- Query: `email`
- Explanation: an `email` that was assigned as a query `?email=` doesn't pass the following field validation rules:
    - email: `required,email,lte=255`

## E009

- Error Name: `SendResetPasswordInvalidEmail`
- Controller: `user`
- Path: `/reset/account`
- Method: `PATCH`
- Status: `400`
- Query: `email`
- Explanation: an `email` that was assigned as a query `?email=` doesn't pass the following field validation rules:
    - email: `required,email,lte=255`

## E010

- Error Name: `UpdatePasswordInvalidBody`
- Controller: `user`
- Path: `/update/password`
- Method: `PATCH`
- Status: `400`
- Content: `application/json`
- Body: `password, token`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - password: `required,gte=5,lte=36`
    - token: `required`

## E011

- Error Name: `UpdatePasswordInvalidToken`
- Controller: `user`
- Path: `/update/password`
- Method: `PATCH`
- Status: `401`
- Content: `application/json`
- Body: `password, token`
- Explanation: the request body contains a `token` that is invalid or expired, a new update password token will need to be regenerated

## E012

- Error Name: `GetEnvironmentInvalidID`
- Controller: `environment`
- Path: `/environment/id/:id`
- Params: `id`
- Controller: `secret`
- Path: `/secrets/search?key=<secret_key>&environmentID=<environmentID>`
- Query: `key, environmentID`
- Method: `GET`
- Status: `400`
- Explanation: the request params or request query doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`
    - environmentID: `required,uuid`

## E013

- Error Name: `GetEnvironmentNonExistentID`
- Controller: `environment`
- Path: `/environment/id/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E014

- Error Name: `GetEnvironmentInvalidName`
- Controller: `environment`
- Path: `/environment/name/?name=<environmentName>&projectID=<projectID>`
- Query: `name, projectID`
- Path: `/environments/search/?name=<environmentName>&projectID=<projectID>`
- Query: `name, projectID`
- Path: `/secrets/projectenvironment/?environment=<environmentName>&project=<projectName>`
- Query: `environment, project`
- Method: `GET`
- Status: `400`
- Explanation: the request query param doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)
    - environment: `required,name,lte=255` (`name` is a custom validation)

## E015

- Error Name: `GetEnvironmentInvalidProjectID`
- Controller: `environment`
- Path: `/environment/name/?name=<environmentName>&projectID=<projectID>`
- Path: `/environments/search/?name=<environmentName>&projectID=<projectID>`
- Method: `GET`
- Status: `400`
- Query: `name, projectID`
- Explanation: the `projectID` query param doesn't pass one or more of the following field validation rules:
    - projectID: `required,uuid`

## E016

- Error Name: `GetEnvironmentNonExistentName`
- Controller: `environment`
- Path: `/environment/name/:name`
- Params: `name`
- Controller: `secret`
- Path: `/secrets/projectenvironment/?environment=<environmentName>&project=<projectName>`
- Query: `environment, project`
- Method: `GET`
- Status: `404`
- Explanation: the request params `name` or request query `environement` contains a value that doesn't match a user created environment

## E017

- Error Name: `CreateEnvironmentInvalidBody`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `400`
- Body: `name, projectID`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)
    - projectID: `required,uuid`

## E018

- Error Name: `CreateEnvironmentInvalidProjectID`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `404`
- Body: `name, projectID`
- Explanation: the request body contains a `projectID` that doesn't match a user created project

## E019

- Error Name: `CreateEnvironmentNameTaken`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `409`
- Params: `name`
- Explanation: the request params contains an environment `name` that is already in use by the user; another 
name should be used instead

## E020

- Error Name: `DeleteEnvironmentInvalidID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E021

- Error Name: `DeleteEnvironmentNonExistentID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E022

- Error Name: `UpdateEnvironmentInvalidBody`
- Controller: `environment`
- Path: `/update/environment`
- Method: `PUT`
- Status: `400`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`
    - projectID: `required,uuid`
    - updatedName: `required,name,lte=255` (`name` is a custom validation)

## E023

- Error Name: `UpdateEnvironmentInvalidProjectID`
- Controller: `environment`
- Path: `/update/environment`
- Method: `PUT`
- Status: `400`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains a `projectID` value that doesn't match any user created projects

## E024

- Error Name: `UpdateEnvironmentNonExistentID`
- Controller: `user`
- Path: `/update/environment`
- Method: `PUT`
- Status: `404`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains an `id` value that doesn't match any user created environment

## E025

- Error Name: `UpdateEnvironmentNameTaken`
- Controller: `user`
- Path: `/update/environment`
- Method: `PUT`
- Status: `409`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains a `name` value that in use by the user; another 
environment name should be used instead

## E026

- Error Name: `GetSecretInvalidID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E027

- Error Name: `GetSecretNonExistentID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created secret

## E028

- Error Name: `GetSecretsByEnvInvalidID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E029

- Error Name: `GetSecretsByEnvNonExistentID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E030

- Error Name: `CreateSecretInvalidBody`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `400`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - environmentIDs: `uuidarray` (`uuidarray` is a custom validation)
    - key: `required,gte=2,lte=255`
    - projectID: `required,uuid`
    - value: `required,lte=5000`

## E031

- Error Name: `CreateSecretNonExistentProject`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `404`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `projectID` value doesn't match any user created projects

## E032

- Error Name: `CreateSecretNonExistentEnv`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `404`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `environmentIDs` value doesn't match any user created environments

## E033

- Error Name: `CreateSecretKeyAlreadyExists`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `409`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `key` value matches a pre-existing key value within one or more user created environments

## E034

- Error Name: `DeleteSecretInvalidID`
- Controller: `secret`
- Path: `/delete/secret/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E035

- Error Name: `DeleteSecretNonExistentID`
- Controller: `secret`
- Path: `/delete/secret/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E036

- Error Name: `UpdateSecretInvalidBody`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `400`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`
    - environmentIDs: `uuidarray` (`uuidarray` is a custom validation)
    - key: `required,gte=2,lte=255`
    - value: `required,lte=5000`

## E037

- Error Name: `UpdateSecretInvalidID`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `404`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body contains an `id` that doesn't match a user created secret

## E038

- Error Name: `UpdateSecretNonExistentEnv`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `404`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body `environmentIDs` value doesn't match any user created environments

## E039

- Error Name: `UpdateSecretKeyAlreadyExists`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `409`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body `key` value matches a pre-existing key value within one or more user created environments

## E040

- Error Name: `GetProjectInvalidID`
- Controller: `project`
- Path: `/project/id/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E041

- Error Name: `GetProjectInvalidID`
- Controller: `project`
- Path: `/project/id/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params `id` value doesn't match any user created projects

## E042

- Error Name: `GetProjectInvalidName`
- Controller: `environment`
- Params: `name`
- Path: `/environments/project/:name`
- Controller: `project`
- Params: `name`
- Path: `/secrets/projectenvironment`
- Controller: `project`
- Path: `/project/name/:name`
- Query: `name`
- Method: `GET`
- Status: `400`
- Explanation: the request params or request query doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)

## E043

- Error Name: `GetProjectNonExistentName`
- Controller: `environment`
- Path: `/environments/project/:name`
- Params: `name`
- Controller: `project`
- Path: `/project/name/:name`
- Params: `name`
- Controller: `secret`
- Path: `/secrets/projectenvironment/?environment=<environmentName>&project=<projectName>`
- Query: `environment, project`
- Method: `GET`
- Status: `404`
- Explanation: the request params `name` or request query `project` value doesn't match any user created projects

## E044

- Error Name: `CreateProjectInvalidName`
- Controller: `project`
- Path: `/create/project/:name`
- Path: `/projects/search/:name`
- Method: `POST`
- Status: `400`
- Params: `name`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)

## E045

- Error Name: `CreateProjectNameTaken`
- Controller: `project`
- Path: `/create/project/:name`
- Method: `POST`
- Status: `409`
- Params: `name`
- Explanation: the request params `name` value matches a project name that already exists

## E046

- Error Name: `DeleteProjectInvalidID`
- Controller: `project`
- Path: `/delete/project/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E047

- Error Name: `DeleteProjectNonExistentID`
- Controller: `project`
- Path: `/delete/project/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created project

## E048

- Error Name: `UpdateProjectInvalidBody`
- Controller: `project`
- Path: `/update/project`
- Method: `PUT`
- Status: `400`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`
    - updatedName: `required,name,lte=255` (`name` is a custom validation)

## E049

- Error Name: `UpdateProjectNonExistentID`
- Controller: `project`
- Path: `/update/project`
- Method: `PUT`
- Status: `404`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body contains an `id` value that doesn't match a user created project

## E050

- Error Name: `UpdateProjectNameTaken`
- Controller: `project`
- Path: `/update/project`
- Method: `PUT`
- Status: `409`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body contains a `name` value that in use by the user; another 
project name should be used instead

## E051

- Error Name: `SearchForSecretsByEnvAndSecretInvalidKey`
- Controller: `secret`
- Path: `/secrets/search?key=<secret_key>&environmentID=<environmentID>`
- Method: `PUT`
- Status: `400`
- Query: `key, environmentID`
- Explanation: the request query doesn't pass one or more of the following field validation rules:
    - key: `required,gte=2,lte=255`

## E052

- Error Name: `CreateProjectOverLimit`
- Controller: `project`
- Path: `/create/project/:name`
- Method: `POST`
- Status: `403`
- Params: `name`
- Explanation: the request is attempting to create a project that goes over the limit of 10 project per account

## E053

- Error Name: `CreateEnvironmentOverLimit`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `403`
- Body: `name, projectID`
- Explanation: the request is attempting to create an environment that goes over the limit of 10 environments per account

## E054

- Error Name: `UpdateDisplayNameMissingName`
- Controller: `user`
- Path: `/update/name`
- Method: `PATCH`
- Status: `400`
- Query: `name`
- Explanation: the request query doesn't pass one or more of the following field validation rules:
    - name: `required,gte=2,lte=64`
