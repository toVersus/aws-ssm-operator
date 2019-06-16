package parameterstore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	errs "github.com/pkg/errors"
	"github.com/toVersus/aws-ssm-operator/pkg/apis/ssm/v1alpha1"
)

// SSMClient preserves AWS client session and SSM client itself
type SSMClient struct {
	s   *session.Session
	ssm *ssm.SSM
}

func newSSMClient(s *session.Session) *SSMClient {
	return &SSMClient{
		s: s,
	}
}

// FetchParameterStoreValue fetches decrypted value from SSM Parameter Store
func (c *SSMClient) FetchParameterStoreValue(name string) (*string, error) {
	if c.s == nil {
		c.s = session.Must(session.NewSession())
	}

	if c.ssm == nil {
		c.ssm = ssm.New(c.s)
	}

	got, err := c.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	return got.Parameter.Value, nil
}

// SSMParameterValueToSecret shapes fetched value so as to store it into K8S Secret
func (c *SSMClient) SSMParameterValueToSecret(ref v1alpha1.ParameterStoreRef) (map[string]string, error) {
	log.Info("fetching value from SSM Parameter Store")
	val, err := c.FetchParameterStoreValue(ref.Name)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, errs.New("fetched value must not be nil")
	}

	return map[string]string{
		"name": *val,
	}, nil
}
