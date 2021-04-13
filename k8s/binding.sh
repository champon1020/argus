#!/bin/zsh

gcloud iam service-accounts add-iam-policy-binding \
       --role roles/iam.workloadIdentityUser \
       --member "serviceAccount:alfheim-argus-269319.svc.id.goog[default/myblog-external-secrets-kubernetes-external-secrets]" \
       myblog-cluster@alfheim-argus-269319.iam.gserviceaccount.com
