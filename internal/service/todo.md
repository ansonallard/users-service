What's next:
- Implement tenants
  - Create an `admin` tenant. On service startup, check database to see if admin user exists in `admin` tenant, if not, prompt user to create admin use
  - Store username and password salted and hashed
- Admin can create tenants, users, and applications
```
POST /tenants 
response: 
{
    tenantId: Ulid
    createdAt: Date
    version: string
}

PUT /tenants/{id}/users/{username}
request:
{
    password: string
}
response 204

POST /tenants/{id}/users/{username}/resetPassword
request: 
{
    oldPassword: string
    newPassword: string
}
response: 204 // requires you to relogin to get access token

POST /tenants/{id}/applications
request:
{
    name: string
    users: string[]
}
response: 
{
    clientId: string
}

GET /tenants/{id}/applications/{id}/users
request {
    nextToken: string | undefined
    maxResults: int // bounded to 1000
}
response:
{
    usernames: string[]
    nextToken: string
}

PUT /tenants/{id}/applications/{id}/users/{username}
response 204

GET /tenants/{id}/oauth/authorize

POST /tenants/{id}/oauth/token
```

Create well known json blob in nginx config at `hostname/.well-known/authorization` containing (can cache this value): 

```json
{
    "tenant_id": "string",
    "oauth_token_endpoint": "string",
    "oauth_authorization_endpoint": "string",
    "jwk_endpoint": "string",
    "client_id": "string"
}
```