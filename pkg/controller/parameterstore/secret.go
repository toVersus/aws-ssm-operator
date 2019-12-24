package parameterstore

import (
	"strings"

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

// FetchParameterStoreValue fetches decrypted values from SSM Parameter Store
func (c *SSMClient) FetchParameterStoreValues(ref v1alpha1.ParameterStoreRef) (map[string]string, error) {
	if c.s == nil {
		c.s = session.Must(session.NewSession())
	}

	if c.ssm == nil {
		c.ssm = ssm.New(c.s)
	}

	if ref.Name != "" {
		log.Info("fetching values from SSM Parameter Store by name: %s", ref.Name)
		got, err := c.ssm.GetParameter(&ssm.GetParameterInput{
			Name:           aws.String(ref.Name),
			WithDecryption: aws.Bool(true),
		})
		if err != nil {
			return nil, err
		}

		return map[string]string{"name": aws.StringValue(got.Parameter.Value)}, nil
	}

	log.Info("fetching values from SSM Parameter Store by path: %s", ref.Path)
	got, err := c.ssm.GetParametersByPath(&ssm.GetParametersByPathInput{
		Path:           aws.String(ref.Path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	dict := make(map[string]string, len(got.Parameters))
	for _, p := range got.Parameters {
		ss := strings.Split(aws.StringValue(p.Name), "/")
		dict[ss[len(ss)-1]] = aws.StringValue(p.Value)
	}

	return dict, nil
}

// SSMParameterValueToSecret shapes fetched value so as to store them into K8S Secret
func (c *SSMClient) SSMParameterValueToSecret(ref v1alpha1.ParameterStoreRef) (map[string]string, error) {
	params, err := c.FetchParameterStoreValues(ref)
	if err != nil {
		return nil, err
	}
	if params == nil {
		return nil, errs.New("fetched value must not be nil")
	}

	return params, nil
}
