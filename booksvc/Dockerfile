# syntax=docker/dockerfile:1

FROM golang:1.17-alpine 
WORKDIR /booksvc
# Insert all dependecies
COPY go.mod ./
COPY go.sum ./
#Installing on our docker immage
RUN go mod download
#Copy all .go files of our project
COPY . .
#Run tests
#RUN CGO_ENABLED=0 go test -v
#building the bim exec and calling docker-booksvc
RUN go build -o ./docker-yale-booksvc
# using the port 8013
EXPOSE 8013
#on start execute comands

CMD [ "./docker-yale-booksvc" ]