package s3p

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"net/url"
	"os"
	"strings"
	"sync"
)

type Location struct {
	Bucket *string
	Key    *string
}

type AccessPoint struct {
	Endpoint    string
	Region      string
	Bucket      string
	Prefix      string
	Credentials *credentials.Credentials
	session     *session.Session
}

var mu sync.Mutex
var registry = map[string]*AccessPoint{}

func Register(ep string, ap AccessPoint) {
	mu.Lock()
	registry[strings.ToLower(ep)] = &ap
	mu.Unlock()
}

func Lookup(ep string) (ap *AccessPoint, ok bool) {
	mu.Lock()
	defer mu.Unlock()
	ep = strings.ToLower(ep)
	ap = registry[ep]
	if ap == nil {
		nx := "s3_" + ep + "_url"
		if ep == "" {
			nx = "s3_url"
		}
		for _, v := range os.Environ() {
			j := strings.Index(v, "=")
			if j > 0 && strings.ToLower(v[:j]) == nx {
				u, err := url.Parse(v[j+1:])
				if err != nil {
					return nil, false
				}
				ap, err = DecodeUrl(u)
				if err != nil {
					return nil, false
				}
				registry[ep] = ap
				break
			}
		}
	}
	return ap, ap != nil
}

func (ap *AccessPoint) Session() (ssn *session.Session, err error) {
	mu.Lock()
	defer mu.Unlock()
	if ap.session != nil {
		return ap.session, nil
	}
	ap.session, err = session.NewSession(&aws.Config{
		Endpoint:    &ap.Endpoint,
		Region:      &ap.Region,
		Credentials: ap.Credentials,
	})
	return ap.session, err
}
