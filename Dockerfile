FROM hashicorp/terraform:0.13.2 as initializer

COPY resources /resources
RUN cd /resources/terraform && terraform init

FROM hashicorp/terraform:0.13.2

ENV M_WORKDIR "/workdir"
ENV M_RESOURCES "/resources"
ENV M_SHARED "/shared"

WORKDIR /workdir
ENTRYPOINT ["make"]

RUN apk add --update --no-cache make=4.3-r0 &&\
    wget https://github.com/mikefarah/yq/releases/download/3.3.4/yq_linux_amd64 -O /usr/bin/yq &&\
    chmod +x /usr/bin/yq &&\
    # Installing helm
    wget https://get.helm.sh/helm-v3.3.4-linux-amd64.tar.gz &&\
    tar -zxvf helm-v3.3.4-linux-amd64.tar.gz &&\
    mv linux-amd64/helm /usr/local/bin/helm &&\
    rm -rf linux-amd64 &&\
    chmod +x /usr/local/bin/helm

ARG ARG_M_VERSION="unknown"
ENV M_VERSION=$ARG_M_VERSION

COPY --from=initializer /resources/ /resources/
COPY workdir /workdir

ARG ARG_HOST_UID=1000
ARG ARG_HOST_GID=1000
RUN chown -R $ARG_HOST_UID:$ARG_HOST_GID \
    $M_WORKDIR \
    $M_RESOURCES

USER $ARG_HOST_UID:$ARG_HOST_GID
# Set HOME to directory with necessary permissions for current user
ENV HOME=$M_WORKDIR

RUN helm repo add stable https://charts.helm.sh/stable
