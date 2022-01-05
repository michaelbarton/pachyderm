#!/bin/bash

set -ex

# shellcheck disable=SC1090
source "$(dirname "$0")/env.sh"

# deploy object storage
kubectl apply -f etc/testing/minio-openshift.yaml

REGISTRY=$(oc get route -n openshift-image-registry | awk 'FNR == 2 {print $2}') &&
PROJECT="pach-test"
PACHD_REPO="${REGISTRY}/${PROJECT}/pachd"
WORKER_REPO="${REGISTRY}/${PROJECT}/worker"

helm install pachyderm etc/helm/pachyderm -f etc/testing/circle/helm-values.yaml \
--set etcd.securityContext.enabled=false --set pachd.securityContext.enabled=false \
--set pachd.image.repository=${PACHD_REPO} --set pachd.worker.image.repository=${WORKER_REPO} \
--set postgresql.enabled=true postgresql.securityContext.enabled=false \
--set postgresql.persistence.enabled=false --set postgresql.containerSecurityContext.enabled=false;

kubectl wait --for=condition=ready pod -l app=pachd --timeout=5m
