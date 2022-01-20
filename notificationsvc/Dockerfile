# syntax=docker/dockerfile:1

FROM golang:1.17-alpine 
WORKDIR /customersvc
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
RUN go build -o ./docker-yale-notificationsvc
# using the port 3000
EXPOSE 3000
#on start execute comands

CMD [ "./docker-yale-notificationsvc" ]