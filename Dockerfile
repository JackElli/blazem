# pulls from golang alpine
FROM golang:alpine as builder

# install git for dependencies
RUN apk update && apk add --no-cache git

# set workpath to anything with root
WORKDIR /app

#copy blazem (has to be local dir so blazem/)
COPY . .

# get dependecies
# and build
RUN go get -d -v
RUN go build -o /blazem 

FROM scratch

COPY --from=builder /blazem /blazem

EXPOSE 3100

ENTRYPOINT [ "/blazem" ]