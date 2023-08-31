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
- Status: `200`
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
- Path: `/environment/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
  - id: `required,uuid`

## E013

- Error Name: `GetEnvironmentNonExistentID`
- Controller: `environment`
- Path: `/environment/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E014

- Error Name: `CreateEnvironmentInvalidName`
- Controller: `environment`
- Path: `/create/environment/:name`
- Method: `POST`
- Status: `400`
- Params: `name`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
  - name: `required,envname,lte=255` (`envname` is a custom validation)

## E015

- Error Name: `CreateEnvironmentNameTaken`
- Controller: `environment`
- Path: `/create/environment/:name`
- Method: `POST`
- Status: `409`
- Params: `name`
- Explanation: the request params contains an environment `name` that already in use by the user; another name should 
be used instead

## E016

- Error Name: `DeleteEnvironmentInvalidID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `400`
- Params: `id`
- Explanation: the request params doesn't pass one or more of the following field validation rules:
  - id: `required,uuid`

## E017

- Error Name: `DeleteEnvironmentNonExistentID`
- Controller: `environment`
- Path: `/delete/environment/:id`
- Method: `DELETE`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E018

- Error Name: `UpdateEnvironmentInvalidBody`
- Controller: `environment`
- Path: `/update/environment`
- Method: `PUT`
- Status: `400`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
  - id: `required,uuid`
  - updatedName: `required,envname,lte=255` (`envname` is a custom validation)

## E019

- Error Name: `UpdateEnvironmentNonExistentID`
- Controller: `user`
- Path: `/update/environment`
- Method: `PUT`
- Status: `404`
- Content: `application/json`
- Body: `id, updatedName`
- Explanation: the request params contains an `id` that doesn't match a user created environment

## E020

- Error Name: `GetSecretInvalidID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
  - id: `required,uuid`

## E021

- Error Name: `GetSecretNonExistentID`
- Controller: `secret`
- Path: `/secret/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created secret

## E022

- Error Name: `GetSecretsByEnvInvalidID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `400`
- Params: `id`
- Explanation: the request body doesn't pass one or more of the following field validation rules:
  - id: `required,uuid`

## E023

- Error Name: `GetSecretsByEnvNonExistentID`
- Controller: `secret`
- Path: `/secrets/:id`
- Method: `GET`
- Status: `404`
- Params: `id`
- Explanation: the request params contains an `id` that doesn't match a user created environment
