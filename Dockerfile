FROM golang:latest

LABEL maintainer="Gábor Görzsöny <gabor@gorzsony.com>"

WORKDIR /app
COPY . .
RUN make

EXPOSE 8080

ENTRYPOINT ["./razlink"]