FROM golang:1.16.4

ENV WORKPATH=/app/src/github.com/jiaqi-yin/go-url-shortener
COPY . ${WORKPATH}
WORKDIR ${WORKPATH}

RUN go build -o go-url-shortener .

EXPOSE 8080

CMD [ "./go-url-shortener" ]