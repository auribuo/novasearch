FROM ubuntu:20.04

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install git curl tar unzip -y

WORKDIR /

RUN curl -sLO https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.20.2.linux-amd64.tar.gz

ADD 'https://api.github.com/repos/auribuo/novasearch/commits?per_page=1' /tmp/last_commit.json
RUN curl -sLO "https://github.com/auribuo/novasearch/archive/main.zip" && unzip main.zip
RUN mv novasearch-main novasearch
WORKDIR /novasearch

ENV PATH="$PATH:/root/go/bin:/usr/local/go/bin"
RUN go build
RUN go install

EXPOSE 8080

ENTRYPOINT ["novasearch"]