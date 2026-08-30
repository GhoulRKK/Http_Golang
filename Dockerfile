FROM golang:1.25-bookworm 
#образ
#директива, файловая система
WORKDIR /app 
#перенести код, файлы
COPY . . 
#установить зависимости

RUN go mod tidy 

#запуск

CMD ["make","service-run"] 