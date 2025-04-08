FROM golang:1.24 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -tags lambda.norpc -o main ./cmd/lambda

FROM public.ecr.aws/lambda/provided:al2023

COPY --from=build /app/main ./main

ENTRYPOINT ["./main"]
