1. Apply `crc.patch` on the repository. Pay attention: `openshift-edo` account's default VPC ID is used (not infoblox one).
2. Build the image: `VERSION=x.y.z ./all -b -c`
3. Deploy the operator: `IMG=quay.io/alebedev/albo:x.y.z make deploy`
4. Create credentials for `openshift-edo` account: `oc create secret generic aws-load-balancer-operator -n aws-load-balancer-operator --from-file=credentials=credentials`
