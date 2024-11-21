# 다이어리 API
## 목록 불러오기
`[GET] /diaries`

### Request
#### Header
```
Authorization: 'Bearer <accessToken>'
```

### Response
#### Status Code
- 200 OK
- 400 BadRequest
- 500 Unknown

#### Body (application/json)
```json
[
  {
    "content": "이건 내용입니다.",
    "createdAt": "2024-11-20T04:36:35+09:00",
    "id": "019347db-9f6a-7ba8-8ddb-f197f13d816b",
    "mood": "exhaust",
    "title": "안녕하세요",
    "updatedAt": "2024-11-20T04:36:35+09:00"
  }
]
```

## 다이어리 생성
`[POST] /diaries`

### Request
#### Header
```
Authorization: 'Bearer <accessToken>'
```

#### Body (application/json)
```json
{
  "title": "안녕하세요",
  "content": "이건 내용입니다.",
  "mood": "exhaust" // "surprise", "angry", "sad", "neutral", "Mad", "cry", "happy", "exhaust"
}
```
- title
  - 255자 이내
- content
  - 1024자 이내
- Mood
  - "surprise", "angry", "sad", "neutral", "mad", "cry", "happy", "exhaust" 중에서 하나만 사용가능

### Response
#### Status Code
- 200 OK
- 400 BadRequest
- 500 Unknown

#### Body (application/json)
```json
{
  "content": "이건 내용입니다.",
  "createdAt": "2024-11-20T04:36:35+09:00",
  "id": "019347db-9f6a-7ba8-8ddb-f197f13d816b",
  "mood": "exhaust",
  "title": "안녕하세요",
  "updatedAt": "2024-11-20T04:36:35+09:00"
}
```

## 다이어리 수정
`[PATCH] /diaries/:id`

### Request
#### Header
```
Authorization: 'Bearer <accessToken>'
```

#### Body (application/json)
```json
{
    "title": "안녕하세요",
    "content": "이건 내용입니다.",
    "mood": "exhaust" // "surprise", "angry", "sad", "neutral", "Mad", "cry", "happy", "exhaust"
}
```
- title
  - 255자 이내
  - nullable
- content
  - 1024자 이내
  - nullable
- Mood
  - "surprise", "angry", "sad", "neutral", "mad", "cry", "happy", "exhaust" 중에서 하나만 사용가능
  - nullable

### Response
#### Status Code
- 200 OK
- 400 BadRequest
- 404 NotFound
- 500 Unknown

#### Body (application/json)
```json
{
  "content": "이건 내용입니다.",
  "createdAt": "2024-11-20T04:36:35+09:00",
  "id": "019347db-9f6a-7ba8-8ddb-f197f13d816b",
  "mood": "exhaust",
  "title": "안녕하세요",
  "updatedAt": "2024-11-20T04:36:35+09:00"
}
```

## 다이어리 삭제
`[DELETE] /diaries/:id`

### Request
#### Header
```
Authorization: 'Bearer <accessToken>'
```

### Response
#### Status Code
- 204 NoContent
- 400 BadRequest
- 404 NotFound
- 500 Unknown

#### Body (application/json)
```json
""
```
