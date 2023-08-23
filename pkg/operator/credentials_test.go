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

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	cco "github.com/openshift/cloud-credential-operator/pkg/apis/cloudcredential/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/openshift/aws-load-balancer-operator/pkg/utils/test"
)

func TestProvisionCredentials(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		scheme  *runtime.Scheme
		// provisionedSecret simulates CCO job
		// which is supposed to create a secret from CredentialsRequest
		provisionedSecret   *corev1.Secret
		expectedCredReqName types.NamespacedName
		compareCredReq      func(*cco.CredentialsRequest, *cco.AWSProviderSpec) error
		expectedContents    string
		errExpected         bool
	}{
		{
			name: "nominal sts",
			envVars: map[string]string{
				"ROLEARN": "arn:aws:iam::123456789012:role/foo",
			},
			scheme: test.Scheme,
			provisionedSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-load-balancer-operator",
					Namespace: "aws-load-balancer-operator",
				},
				Data: map[string][]byte{
					"credentials": []byte("oksts"),
				},
			},
			expectedCredReqName: types.NamespacedName{
				Namespace: "openshift-cloud-credential-operator",
				Name:      "aws-load-balancer-operator",
			},
			compareCredReq: func(credReq *cco.CredentialsRequest, providerSpec *cco.AWSProviderSpec) error {
				if providerSpec.STSIAMRoleARN != "arn:aws:iam::123456789012:role/foo" {
					return fmt.Errorf("got unexpected role arn: %q", providerSpec.STSIAMRoleARN)
				}
				if credReq.Spec.CloudTokenPath != "/var/run/secrets/openshift/serviceaccount/token" {
					return fmt.Errorf("got unexpected token path: %q", credReq.Spec.CloudTokenPath)
				}
				return nil
			},
			expectedContents: "oksts",
		},
		{
			name:   "nominal non sts",
			scheme: test.Scheme,
			provisionedSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-load-balancer-operator",
					Namespace: "aws-load-balancer-operator",
				},
				Data: map[string][]byte{
					"credentials": []byte("oknonsts"),
				},
			},
			expectedCredReqName: types.NamespacedName{
				Namespace: "openshift-cloud-credential-operator",
				Name:      "aws-load-balancer-operator",
			},
			compareCredReq: func(credReq *cco.CredentialsRequest, providerSpec *cco.AWSProviderSpec) error {
				if providerSpec.STSIAMRoleARN != "" {
					return fmt.Errorf("expected role arn to be unset but got %q", providerSpec.STSIAMRoleARN)
				}
				if credReq.Spec.CloudTokenPath != "" {
					return fmt.Errorf("expected token path to be unset but got %q", credReq.Spec.CloudTokenPath)
				}
				return nil
			},
			expectedContents: "oknonsts",
		},
		{
			name:   "invalid role arn",
			scheme: test.Scheme,
			envVars: map[string]string{
				"ROLEARN": "arn:aws:iam:role/foo",
			},
			provisionedSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-load-balancer-operator",
					Namespace: "aws-load-balancer-operator",
				},
				Data: map[string][]byte{
					"credentials": []byte("ok"),
				},
			},
			errExpected: true,
		},
		{
			name:   "credentialsrequest creation failed",
			scheme: test.BasicScheme,
			provisionedSecret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "aws-load-balancer-operator",
					Namespace: "aws-load-balancer-operator",
				},
				Data: map[string][]byte{
					"credentials": []byte("oknonsts"),
				},
			},
			errExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envVars {
				err := os.Setenv(k, v)
				if err != nil {
					t.Fatalf("failed to set %q environment variable: %v", k, err)
				}
				defer os.Unsetenv(k)
			}

			cli := fake.NewClientBuilder().WithScheme(tc.scheme).WithObjects(tc.provisionedSecret).Build()

			gotFilename, err := ProvisionCredentials(context.Background(), cli, tc.provisionedSecret.Namespace)
			if err != nil {
				if !tc.errExpected {
					t.Fatalf("got unexpected error: %v", err)
				}
				return
			} else if tc.errExpected {
				t.Fatalf("error expected but not received")
			}

			if tc.scheme != test.BasicScheme {
				credReq := &cco.CredentialsRequest{}
				err = cli.Get(context.Background(), tc.expectedCredReqName, credReq)
				if err != nil {
					t.Fatalf("failed to get credentials request %v: %v", tc.expectedCredReqName, err)
				}

				providerSpec := &cco.AWSProviderSpec{}
				err = cco.Codec.DecodeProviderSpec(credReq.Spec.ProviderSpec, providerSpec)
				if err != nil {
					t.Fatalf("failed to decode credentials request's aws provider spec: %v", err)
				}

				if err := tc.compareCredReq(credReq, providerSpec); err != nil {
					t.Fatalf("credentials request comparison failed: %v", err)
				}
			}

			gotContents, err := os.ReadFile(gotFilename)
			if err != nil {
				t.Fatalf("failed to read generated file: %v", err)
			}
			if string(gotContents) != tc.expectedContents {
				t.Fatalf("expected contents %q but got %q", tc.expectedContents, string(gotContents))
			}
		})
	}
}

func TestWaitForSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "test",
		},
	}

	// secret exists
	cli := fake.NewClientBuilder().WithScheme(test.Scheme).WithObjects(secret).Build()
	resCh, errCh := make(chan *corev1.Secret), make(chan error)
	go func(resCh chan *corev1.Secret, errCh chan error) {
		s, err := waitForSecret(context.Background(), cli, secret.Namespace, secret.Name, 10*time.Second, 10*time.Millisecond)
		resCh <- s
		errCh <- err
	}(resCh, errCh)
wait:
	for {
		select {
		case gotSecret := <-resCh:
			sameNs := gotSecret.Namespace == secret.Namespace
			sameName := gotSecret.Name == secret.Name
			if !sameNs || !sameName {
				t.Fatalf("got unexpected secret: %v", gotSecret)
			}
			break wait
		case err := <-errCh:
			t.Fatalf("got unexpected error: %v", err)
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out")
		}
	}

	// secret doesn't exist - timeout
	emptyCli := fake.NewClientBuilder().WithScheme(test.Scheme).WithObjects().Build()
	go func(resCh chan *corev1.Secret, errCh chan error) {
		s, err := waitForSecret(context.Background(), emptyCli, secret.Namespace, secret.Name, 1*time.Second, 10*time.Millisecond)
		resCh <- s
		errCh <- err
	}(resCh, errCh)
waitTimeout:
	for {
		select {
		case gotSecret := <-resCh:
			t.Fatalf("got unexpected secret: %v", gotSecret)
		case <-errCh:
			break waitTimeout
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out")
		}
	}
}

func TestCredentialsFileFromSecret(t *testing.T) {
	tests := []struct {
		name             string
		secret           *corev1.Secret
		pattern          string
		expectedPrefix   string
		expectedContents string
		errExpected      bool
	}{
		{
			name: "nominal",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"credentials": []byte("ok"),
				},
			},
			pattern:          "test-",
			expectedPrefix:   "/tmp/test-",
			expectedContents: "ok",
		},
		{
			name: "wrong data key",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"wrongkey": []byte("ok"),
				},
			},
			pattern:     "test-",
			errExpected: true,
		},
		{
			name: "invalid pattern",
			secret: &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "test",
				},
				Data: map[string][]byte{
					"credentials": []byte("ok"),
				},
			},
			pattern:     "test-//",
			errExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := credentialsFileFromSecret(tc.secret, tc.pattern)
			if err != nil {
				if !tc.errExpected {
					t.Fatalf("got unexpected error: %v", err)
				}
				return
			} else if tc.errExpected {
				t.Fatalf("error expected but not received")
			}
			if !strings.HasPrefix(got, tc.expectedPrefix) {
				t.Fatalf("expected %q to have %q prefix", got, tc.expectedPrefix)
			}
			gotContents, err := os.ReadFile(got)
			if err != nil {
				t.Fatalf("error while reading generated file: %v", err)
			}
			if string(gotContents) != tc.expectedContents {
				t.Fatalf("expected contents %q but got %q", tc.expectedContents, string(gotContents))
			}
		})
	}
}
