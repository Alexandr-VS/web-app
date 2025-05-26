FROM alexvscode/webapp_go

WORKDIR /web-app/cmd/myapp/
RUN git pull
RUN go mod tidy

CMD ["go", "run", "main.go"]
