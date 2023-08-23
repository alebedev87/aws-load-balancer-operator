/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

// NOTE: this file is influenced by the new approach of how to do STS described in
// https://docs.google.com/document/d/1iFNpyycby_rOY1wUew-yl3uPWlE00krTgr9XHDZOTNo.

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	operatorCredentialsRequestName = "aws-load-balancer-operator"
	operatorCredentialsSecretName  = "aws-load-balancer-operator"
	operatorServiceAccountName     = "aws-load-balancer-operator-controller-manager"
	roleARNEnvVar                  = "ROLEARN"
	webIdentityTokenPath           = "/var/run/secrets/openshift/serviceaccount/token"
	credentialsKey                 = "credentials"
	waitForSecretTimeout           = 5 * time.Minute
	waitForSecretPollInterval      = 5 * time.Second
)

// ProvisionCredentials provisions cloud credentials secret in the given namespace
// with IAM policies required by the operator. The credentials data are put
// into a file which can be used to setup AWS SDK client.
func ProvisionCredentials(ctx context.Context, client client.Client, secretNamespace string) (string, error) {
	roleARN := os.Getenv(roleARNEnvVar)
	if roleARN != "" && !arn.IsARN(roleARN) {
		return "", fmt.Errorf("provided role arn is invalid: %q", roleARN)
	}

	credReq := buildCredentialsRequest(secretNamespace, operatorCredentialsSecretName, roleARN)

	// create credentials request resource
	if err := client.Create(ctx, credReq); err != nil {
		if !errors.IsAlreadyExists(err) {
			return "", err
		}
	}

	// wait till the credentials secret is provisioned by CCO
	secret, err := waitForSecret(ctx, client, secretNamespace, operatorCredentialsSecretName, waitForSecretTimeout, waitForSecretPollInterval)
	if err != nil {
		return "", err
	}

	// create credentials file with data taken from the provisioned secret
	credFileName, err := credentialsFileFromSecret(secret, "albo-aws-shared-credentials-")
	if err != nil {
		return "", err
	}

	return credFileName, nil
}

// buildCredentialsRequest returns CredentialsRequest object with IAM policies
// required by this operator. STS IAM role is set if the given role ARN is not empty.
func buildCredentialsRequest(secretNamespace, secretName, roleARN string) *cco.CredentialsRequest {
	providerSpecIn := cco.AWSProviderSpec{
		StatementEntries: GetIAMPolicy().Statement,
	}
	if roleARN != "" {
		providerSpecIn.STSIAMRoleARN = roleARN
	}

	providerSpec, _ := cco.Codec.EncodeProviderSpec(providerSpecIn.DeepCopyObject())

	credReq := &cco.CredentialsRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      operatorCredentialsRequestName,
			Namespace: "openshift-cloud-credential-operator",
		},
		Spec: cco.CredentialsRequestSpec{
			ProviderSpec: providerSpec,
			SecretRef: corev1.ObjectReference{
				Name:      secretName,
				Namespace: secretNamespace,
			},
			ServiceAccountNames: []string{operatorServiceAccountName},
		},
	}

	if roleARN != "" {
		credReq.Spec.CloudTokenPath = webIdentityTokenPath
	}

	return credReq
}

// waitForSecret waits until the secret with the given name appears in the given namespace.
// It returns the secret object and an error if the timeout is exceeded.
func waitForSecret(ctx context.Context, client client.Client, namespace, name string, timeout, pollInterval time.Duration) (*corev1.Secret, error) {
	timeoutCh := time.After(timeout)
	ticker := time.NewTicker(pollInterval)

	for {
		select {
		case <-timeoutCh:
			return nil, fmt.Errorf("timed out waiting for operator credentials secret %q in namespace %q", name, namespace)
		case <-ticker.C:
			secret := &corev1.Secret{}
			err := client.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, secret)
			if err != nil {
				if errors.IsNotFound(err) {
					continue
				} else {
					return nil, err
				}
			} else {
				return secret, nil
			}
		}
	}
}

// credentialsFileFromSecret creates a file on a temporary file system with data from the given secret.
// It returns the full path of the created file and an error
func credentialsFileFromSecret(secret *corev1.Secret, pattern string) (string, error) {
	if len(secret.Data[credentialsKey]) == 0 {
		return "", fmt.Errorf("failed to to find credentials in secret")
	}

	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create credentials file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(secret.Data[credentialsKey]); err != nil {
		return "", fmt.Errorf("failed to write credentials to %q: %w", f.Name(), err)
	}

	return f.Name(), nil
}
