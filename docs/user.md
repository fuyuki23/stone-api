# 유저 API
## 로그인
`[POST] /users/login`

### Request
#### Body (application/json)
```json
{
  "email": "test@test.com",
  "password": "test1111"
}
```

### Response
#### Status Code
- 200 OK
- 400 BadRequest, Invalid Credentials
- 500 Unknown

#### Body (application/json)
```json
{
  "user": {
    "id": "uuid",
    "email": "email@email",
    "name": null, // or "Test"
    "createdAt": "2024-11-18T07:40:57+09:00"
  },
  "tokens": {
    "accessToken": "~~~~", // 15분
    "refreshToken": "~~~~" // 15분 (수정 예정)
  }
}
```

## 회원가입
`[POST] /users/register`

### Request
#### Body (application/json)
```json
{
  "email": "email@email",
  "password": "test1111",
  "name": "Test"
}
```
- email
    - Email 포멧
    - 100자 이내
- password
    - 6자 이상
    - 32자 이내
- name
    - nullable
    - 50자 이내
    - '' 빈 스트링이면 null로 치환

### Response
#### Status Code
- 201 Created
- 400 BadRequest, Invalid Credentials
- 409 UserAlreadyExists
- 500 Unknown

#### Body (application/json)
```json
"ok"
```

## 내 정보 조회
`[GET] /users/me`

### Request
#### Header
```
Authorization: 'Bearer <accessToken>'
```

### Response
#### Status Code
- 200 OK
- 401 Unauthorized

#### Body (application/json)
```json
{
  "id": "uuid",
  "email": "email@email",
  "name": null, // or "Test",
  "createdAt": "2024-11-18T07:40:57+09:00"
}
```

## Refresh Token
`[POST] /users/refresh` a추가 예정...
