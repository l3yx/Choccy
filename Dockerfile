FROM --platform=linux/amd64 ubuntu:20.04

WORKDIR /root

RUN apt-get update && apt-get install -y wget unzip && \
    wget https://github.com/github/codeql/archive/refs/tags/codeql-cli/v2.19.2.tar.gz -O ql.tar.gz && \
    wget https://github.com/github/codeql-cli-binaries/releases/download/v2.19.2/codeql-linux64.zip -O codeql.zip && \
    tar -xf ql.tar.gz && \
    unzip -q codeql.zip && \
    rm -f codeql.zip ql.tar.gz && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

ENV PATH=/root/codeql:$PATH

COPY ./choccy_linux_amd64 ./choccy

RUN chmod +x choccy

# choccy_version=tmp
# docker build . -t l3yx/choccy:${choccy_version}
# docker run --rm -e TZ=Asia/Shanghai -p 22222:80 l3yx/choccy:${choccy_version} ./choccy
# docker tag l3yx/choccy:${choccy_version} l3yx/choccy:latest
# docker push l3yx/choccy:${choccy_version}
# docker push l3yx/choccy:latest
