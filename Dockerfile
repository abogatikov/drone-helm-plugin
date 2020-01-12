ARG golang_version=1.13.6
ARG docker_helm_version=v0.0.1
ARG golangci_version=v1.22.2

FROM golang:${golang_version}-alpine as golang-base
RUN apk add --no-cache make gcc libc-dev git && \
    wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${golangci_version}

FROM golang-base as app
COPY . /workdir
WORKDIR /workdir
ARG release=1.0.0
ENV RELEASE ${release}
RUN make build -e

FROM docker.pkg.github.com/abogatikov/docker-helm/docker-helm:${docker_helm_version} as plugin
LABEL maintainer.name="Alex Bogatikov"
LABEL maintainer.email="a.bogatikov@devalexb.com"
COPY ./env /opt/
COPY --from=app /workdir/build/app /bin/app
ENTRYPOINT ["/bin/app"]