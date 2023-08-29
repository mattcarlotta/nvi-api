# NVI-API Error Codes

Click here for [validation rules](https://github.com/go-playground/validator#baked-in-validations)

## E001

- Controller: `user`
- Path: `/register`
- Content: `application/json`
- Reason: request body is missing a valid `name`, `email` and/or `password`

## E002

- Controller: `user`
- Path: `/register`
- Content: `application/json`
- Reason: request body doesn't pass one or more of the following validation rules:
  - name: `required,gte=2,lte=255`
  - email: `required,email,lte=100`
  - password: `required,gte=5,lte=36`

## E003

- Controller: `user`
- Path: `/register`
- Content: `application/json`
- Reason: request body contains an email that's already in use
