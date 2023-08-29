# NVI-API Error Codes

Click here for [validation rules](https://github.com/go-playground/validator#baked-in-validations)

## E001 - RegisterEmptyBody

- Controller: `user`
- Path: `/register`
- Method: `POST`
- Content: `application/json`
- Reason: request body is missing a valid `name`, `email` and/or `password`

## E002 - RegisterInvalidBody

- Controller: `user`
- Path: `/register`
- Method: `POST`
- Content: `application/json`
- Body: `name, email, password`
- Reason: request body doesn't pass one or more of the following validation rules:
  - name: `required,gte=2,lte=255`
  - email: `required,email,lte=100`
  - password: `required,gte=5,lte=36`

## E003 - RegisterEmailTaken

- Controller: `user`
- Path: `/register`
- Method: `POST`
- Content: `application/json`
- Body: `name, email, password`
- Reason: request body contains an email that's already in use


## E004 - LoginEmptyBody

- Controller: `user`
- Path: `/login`
- Method: `POST`
- Content: `application/json`
- Body: `email, password`
- Reason: request body is missing a valid `email` and/or `password`


## E005 - LoginInvalidBody

- Controller: `user`
- Path: `/login`
- Method: `POST`
- Content: `application/json`
- Body: `email, password`
- Reason: request body doesn't pass one or more of the following validation rules:
  - email: `required,email,lte=100`
  - password: `required,gte=5,lte=36`


## E006 - LoginUnregisteredEmail

- Controller: `user`
- Path: `/login`
- Method: `POST`
- Content: `application/json`
- Body: `email, password`
- Reason: request body contains an unregistered `email`

## E007 - LoginInvalidPassword

- Controller: `user`
- Path: `/login`
- Method: `POST`
- Content: `application/json`
- Body: `email, password`
- Reason: request body contains an invalid `password` for the provided `email`

## E008 - LoginAccountNotVerified

- Controller: `user`
- Path: `/login`
- Method: `POST`
- Content: `application/json`
- Body: `email, password`
- Reason: request body contains an `email` that hasn't been verified yet

