FROM golang:latest AS BUILD 
WORKDIR /root/
COPY main.go .
# RUN go build -o flightpath main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o flightpath main.go


FROM alpine:latest
WORKDIR /root/
COPY --from=BUILD /root/flightpath .
CMD ["./flightpath"]