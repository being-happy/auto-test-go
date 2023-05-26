#
# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .

RUN go build -o /app/output/auto-test-engine  /app/cmd/engine/main.go

FROM debian:buster-20221205

WORKDIR /app

COPY --from=build /app/output/auto-test-engine /app/auto-test-engine

COPY --from=build /app/lua /app/lua
COPY --from=build /app/conf /app/conf

RUN  sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && sed -i s@/deb.debian.org/@/mirrors.aliyun.com/@g /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl telnet tcpdump net-tools vim inetutils-ping \
    && mkdir /app/database

ENV LUA_PATH /app

EXPOSE 8090

CMD ["/app/auto-test-engine"]
