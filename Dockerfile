FROM golang:1.25-bookworm 

WORKDIR /app 

COPY . . 


RUN go mod tidy 


CMD ["make","service-run"] 