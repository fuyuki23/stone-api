# Stone API
돌맹 다이어리 API 서버입니다.

### Requirements
- Go 1.22
- Docker + Docker Compose
- Make
- [air](https://github.com/air-verse/air)

## Development
```shell
make db-up
make watch
```

## Preview

### Requirements
- Docker + Docker Compose

### Start
인증 토큰 생성용 키를 만들어야 합니다.
```shell
openssl genrsa 2048 | awk '{$1=$1;print}' | awk 1 ORS='\\n' | sed 's/..$//' | pbcopy
```
를 실행하여 키를 클립보드에 복사합니다.

`config.toml.tmpl`를 복사하여 `config.toml`을 만들고, `[server.jwt]privateKey`에 복사한 키를 붙여넣습니다. 
`database` 설정은 적절하게 변경해줍니다. `preview` 모드로 실행할 때는 [docker-compose.preview.yaml](./docker/docker-compose.preview.yaml)를 참고하여 작성해줍니다. (Preview 모드는 MySQL을 기본 포트로 사용합니다. 포트: 3306)

```shell
make preview-start
```

### Stop
```shell
make preview-stop
```