#!/bin/bash

# standard bash error handling
set -o errexit;
set -o pipefail;
set -o nounset;
# debug commands
set -x;

# working dir to install binaries etc, cleaned up on exit
BIN_DIR="$(mktemp -d)"
# kind binary will be here
KIND="${BIN_DIR}/kind"
# kubectl binary will be here
KUBECTL="${BIN_DIR}/kubectl"

# cleanup on exit (useful for running locally)
cleanup() {
    "${KIND}" delete cluster || true
    rm -rf "${BIN_DIR}"
}
trap cleanup EXIT

# util to install the latest kind version into ${BIN_DIR}
install_latest_kind() {
    # clone kind into a tempdir within BIN_DIR
    local tmp_dir
    tmp_dir="$(TMPDIR="${BIN_DIR}" mktemp -d "${BIN_DIR}/kind-source.XXXXX")"
    cd "${tmp_dir}" || exit
    git clone https://github.com/kubernetes-sigs/kind && cd ./kind
    make install INSTALL_DIR="${BIN_DIR}"
}

# util to install a released kind version into ${BIN_DIR}
install_kind_release() {
    # VERSION="v0.10.0"
 
    curl -Lo ./kind "https://kind.sigs.k8s.io/dl/v0.10.0/kind-$(uname)-amd64"
    chmod +x ./kind
    mv ./kind "${KIND}"
}

install_kubectl() {
    curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.15.0/bin/linux/amd64/kubectl
    chmod +x ./kubectl
    sudo mv ./kubectl "${KUBECTL}"
}


main() {
   
    # TODO: invoke your tests here
    # Kubernetes <1.16
    kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.2/cert-manager-legacy.yaml
    kubectl create namespace kube-diagnoser
    kubectl apply -f config/deploy
    until kubectl get pods -n kube-diagnoser | grep master | grep 1/1; do sleep 1; done
    until kubectl get pods -n kube-diagnoser | grep agent | grep 1/1; do sleep 1; done

    # teardown will happen automatically on exit
}

main

