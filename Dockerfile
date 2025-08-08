#
# Контейнер сборки
#
FROM golang:1.24 as builder


ENV CGO_ENABLED=0

COPY . /go/src/github.com/web-rabis/nlrk-elastic-load
WORKDIR /go/src/github.com/web-rabis/nlrk-elastic-load
RUN \
    version=git describe --abbrev=6 --always --tag; \
    echo "version=$version" && \
    cd cmd/ && \
    go build -a -tags nlrk-elastic-load -installsuffix nlrk-elastic-load -ldflags "-X main.version=${version} -s -w" -o /go/bin/nlrk-elastic-load -mod vendor

#
# Контейнер для получения актуальных SSL/TLS сертификатов
#
FROM alpine:3.16 as alpine
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
RUN addgroup -S nlrk-elastic-load && adduser -S nlrk-elastic-load -G nlrk-elastic-load

# копируем документацию
#RUN mkdir -p /usr/share/nlrk-elastic-load
#COPY --from=builder /go/src/github.com/web-rabis/nlrk-elastic-load/api /usr/share/api
#RUN chown -R nlrk-elastic-load:nlrk-elastic-load /usr/share/nlrk-elastic-load

ENTRYPOINT [ "/bin/nlrk-elastic-load" ]

#
# Контейнер рантайма
#
FROM scratch
COPY --from=builder /go/bin/nlrk-elastic-load /bin/nlrk-elastic-load

# копируем сертификаты из alpine
COPY --from=alpine /etc/ssl/certs /etc/ssl/certs

## копируем документацию
#COPY --from=alpine /usr/share/nlrk-elastic-load /usr/share/nlrk-elastic-load

# копируем пользователя и группу из alpine
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group

USER nlrk-elastic-load

ENTRYPOINT ["/bin/nlrk-elastic-load"]



