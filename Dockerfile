FROM golang:1.23.4


WORKDIR /app
COPY . .

RUN go mod tidy

RUN rm -f plugins/filemanager/filemanager.so
RUN rm -f plugins/notes/notes.so

RUN go build -buildmode=plugin -o app/plugins/filemanager/filemanager.so plugins/filemanager/filemanager.go
RUN go build -buildmode=plugin -o app/plugins/notes/notes.so plugins/notes/notes.go

RUN go build main .

EXPOSE 8080

CMD ["./main"]
