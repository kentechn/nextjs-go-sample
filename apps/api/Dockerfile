# syntax=docker/dockerfile:1.19
# ビルドコンテキストは apps/api。OpenAPI 仕様は生成済みコードに埋め込まれるため実行時には不要。
ARG GO_VERSION=1.26.6

# 依存解決のみを行う層。go.mod/go.sum が変わらない限りキャッシュが再利用される。
FROM golang:${GO_VERSION}-alpine AS deps
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local
COPY go.mod go.sum ./
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    go mod download

FROM deps AS build
COPY . .
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ローカル開発用: air によるホットリロード。ソースは compose でバインドマウントする。
FROM deps AS dev
ENV PORT=8080
RUN --mount=type=cache,id=go-mod,target=/go/pkg/mod \
    --mount=type=cache,id=go-build,target=/root/.cache/go-build \
    go build -o /usr/local/bin/air github.com/air-verse/air
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

FROM gcr.io/distroless/static-debian12:nonroot AS runtime
ENV PORT=8080
COPY --from=build /out/api /api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]
