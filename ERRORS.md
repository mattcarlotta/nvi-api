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
    - name: `required,gte=2,lte=255`
    - email: `required,email,lte=100`
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
    - email: `required,email,lte=100`
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
    - email: `required,email,lte=100`

## E009

- Error Name: `SendResetPasswordInvalidEmail`
- Controller: `user`
- Path: `/reset/account`
- Method: `PATCH`
- Status: `400`
- Query: `email`
- Explanation: an `email` that was assigned as a query `?email=` doesn't pass the following field validation rules:
    - email: `required,email,lte=100`

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
- Explanation: the request body contains a `token` that is invalid or expired, a new update password token will need

## E012

- Error Name: `GetEnvironmentInvalidToken`
- Controller: `environment`
- Path: `/environment/id/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

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
- Path: `/environment/name/:name`
- Method: `GET`
- Status: `400`
- Params: `name`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)

## E015

- Error Name: `GetEnvironmentNonExistentName`
- Controller: `environment`
- Path: `/environment/name/:name`
- Method: `GET`
- Status: `404`
- Params: `name`
- Explanation: the request params contains an `name` that doesn't match a user created environment

## E016

- Error Name: `CreateEnvironmentInvalidBody`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `400`
- Body: `name, projectID`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)
    - projectID: `required,uuid`

## E017

- Error Name: `CreateEnvironmentInvalidProjectID`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `404`
- Body: `name, projectID`
- Explanation: the request body contains a `projectID` that doesn't match a user created project

## E018

- Error Name: `CreateEnvironmentNameTaken`
- Controller: `environment`
- Path: `/create/environment`
- Method: `POST`
- Status: `409`
- Params: `name`
- Explanation: the request params contains an environment `name` that is already in use by the user; another 
name should be used instead

## E019

- Error Name: `DeleteEnvironmentInvalidID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E020

- Error Name: `DeleteEnvironmentNonExistentID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E021

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

## E022

- Error Name: `UpdateEnvironmentInvalidProjectID`
- Controller: `environment`
- Path: `/update/environment`
- Method: `PUT`
- Status: `400`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains a `projectID` value that doesn't match any user created projects

## E023

- Error Name: `UpdateEnvironmentNonExistentID`
- Controller: `user`
- Path: `/update/environment`
- Method: `PUT`
- Status: `404`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains an `id` value that doesn't match any user created environment

## E024

- Error Name: `UpdateEnvironmentNameTaken`
- Controller: `user`
- Path: `/update/environment`
- Method: `PUT`
- Status: `409`
- Content: `application/json`
- Body: `id, projectID, updatedName`
- Explanation: the request body contains a `name` value that in use by the user; another 
environment name should be used instead

## E025

- Error Name: `GetSecretInvalidID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E026

- Error Name: `GetSecretNonExistentID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created secret

## E027

- Error Name: `GetSecretsByEnvInvalidID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E028

- Error Name: `GetSecretsByEnvNonExistentID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E029

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

## E030

- Error Name: `CreateSecretNonExistentProject`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `404`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `projectID` value doesn't match any user created projects

## E031

- Error Name: `CreateSecretNonExistentEnv`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `404`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `environmentIDs` value doesn't match any user created environments

## E032

- Error Name: `CreateSecretKeyAlreadyExists`
- Controller: `secret`
- Path: `/create/secret`
- Method: `POST`
- Status: `409`
- Body: `environmentIDs, key, projectID, value`
- Explanation: the request body `key` value matches a pre-existing key value within one or more user created environments

## E033

- Error Name: `DeleteSecretInvalidID`
- Controller: `secret`
- Path: `/delete/secret/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E034

- Error Name: `DeleteSecretNonExistentID`
- Controller: `secret`
- Path: `/delete/secret/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E035

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

## E036

- Error Name: `UpdateSecretInvalidID`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `404`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body contains an `id` that doesn't match a user created secret

## E037

- Error Name: `UpdateSecretNonExistentEnv`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `404`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body `environmentIDs` value doesn't match any user created environments

## E038

- Error Name: `UpdateSecretKeyAlreadyExists`
- Controller: `secret`
- Path: `/update/secret`
- Method: `PUT`
- Status: `409`
- Body: `id, environmentIDs, key, value`
- Explanation: the request body `key` value matches a pre-existing key value within one or more user created environments

## E039

- Error Name: `GetProjectInvalidID`
- Controller: `project`
- Path: `/project/id/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E040

- Error Name: `GetProjectInvalidID`
- Controller: `project`
- Path: `/project/id/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params `id` value doesn't match any user created projects

## E041

- Error Name: `GetProjectInvalidName`
- Controller: `project`
- Path: `/project/name/:name`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)

## E042

- Error Name: `GetProjectNonExistentName`
- Controller: `project`
- Path: `/project/name/:name`
- Method: `GET`
- Status: `404`
- Params: `name`
- Explanation: the request params `name` value doesn't match any user created projects

## E043

- Error Name: `CreateProjectInvalidName`
- Controller: `project`
- Path: `/create/project/:name`
- Method: `POST`
- Status: `400`
- Params: `name`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - name: `required,name,lte=255` (`name` is a custom validation)

## E044

- Error Name: `CreateProjectNameTaken`
- Controller: `project`
- Path: `/create/project/:name`
- Method: `POST`
- Status: `409`
- Params: `name`
- Explanation: the request params `name` value matches a project name that already exists

## E045

- Error Name: `DeleteProjectInvalidID`
- Controller: `project`
- Path: `/delete/project/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
    - id: `required,uuid`

## E046

- Error Name: `DeleteProjectNonExistentID`
- Controller: `project`
- Path: `/delete/project/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created project

## E047

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

## E048

- Error Name: `UpdateProjectNonExistentID`
- Controller: `project`
- Path: `/update/project`
- Method: `PUT`
- Status: `404`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body contains an `id` value that doesn't match a user created project

## E049

- Error Name: `UpdateProjectNameTaken`
- Controller: `project`
- Path: `/update/project`
- Method: `PUT`
- Status: `409`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body contains a `name` value that in use by the user; another 
project name should be used instead
