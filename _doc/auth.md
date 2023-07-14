# Authentication

## Login

1. Client submit login credentials to `/auth/login`
2. This app retrieves user account based on credentials.
3. Use user's account data to generate a JWT token.
4. The JWT token, together with user account fields, are returned to client in JSON format. Example data:

```json
{
    "id": 75,
    "userName": "weiguo.ni",
    "email": "weiguo.ni@ftchinese.com",
    "displayName": "Ni Weiguo",
    "expiresAt": 1021697161632000,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjc1LCJuYW1lIjoid2VpZ3VvLm5pIiwiZXhwIjoxMDIxNjk3MTYxNjMyMDAwLCJpYXQiOjE2ODkzMTQwOTAsImlzcyI6ImNvbS5mdGNoaW5lc2Uuc3VwZXJ5YXJkIn0.wr6bHos3eAY0OxKv3TA1nAQFrHlkRw-ftUXs-2awoqE"
}
```

5. Client then stores the returned data in local storage.
6. Subsequent request from client should contain headers `Authorization: Bearer <jwt token>` where `jwt token` is the value of `token` field.
