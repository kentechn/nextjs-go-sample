# buildx bake 用の定義。`task docker:build` から利用する。
# CI では `--set *.cache-from=type=gha --set *.cache-to=type=gha,mode=max` で
# キャッシュのバックエンドだけを差し替える。

variable "TAG" {
  default = "dev"
}

variable "CACHE_DIR" {
  default = ".buildx-cache"
}

group "default" {
  targets = ["api", "web"]
}

target "_common" {
  platforms = ["linux/amd64"]
  # docker-container ドライバでもローカルの docker images に取り込む。
  output = ["type=docker"]
}

target "api" {
  inherits   = ["_common"]
  context    = "apps/api"
  dockerfile = "Dockerfile"
  target     = "runtime"
  tags       = ["nextjs-go-sample/api:${TAG}"]
  cache-from = ["type=local,src=${CACHE_DIR}/api"]
  cache-to   = ["type=local,dest=${CACHE_DIR}/api,mode=max"]
}

target "web" {
  inherits   = ["_common"]
  context    = "."
  dockerfile = "apps/web/Dockerfile"
  target     = "runtime"
  tags       = ["nextjs-go-sample/web:${TAG}"]
  cache-from = ["type=local,src=${CACHE_DIR}/web"]
  cache-to   = ["type=local,dest=${CACHE_DIR}/web,mode=max"]
}
