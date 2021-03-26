package gcp

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"sudachen.xyz/pkg/go-data/errors"
	"sudachen.xyz/pkg/go-data/fu"
	"sync"
)

type Location struct {
	Bucket string
	Key    string
}

type AccessPoint struct {
	Bucket      string
	Prefix      string
	Credentials *google.Credentials
	service     *storage.Service
}

var mu sync.Mutex
var registry = map[string]*AccessPoint{}

func Register(ep string, ap *AccessPoint) {
	mu.Lock()
	registry[strings.ToLower(ep)] = ap
	mu.Unlock()
}

func Lookup(gsurl string) (ap *AccessPoint, loc Location, err error) {
	u, err := url.Parse(gsurl)
	if err != nil {
		return
	}
	if len(u.Host) == 0 {
		err = errors.New("bad gsurl (empty hostname) `" + gsurl + "`")
		return
	}
	if u.Host[0] != '$' {
		err = errors.New("bad gsurl (hostname must start with '$') `" + gsurl + "`")
		return
	}
	mu.Lock()
	defer mu.Unlock()
	ep := strings.ToLower(u.Host[1:])
	ap = registry[strings.ToLower(ep)]
	if ap == nil {
		if ep == "" {
			ap = &AccessPoint{}
			if ap.Credentials, err = google.FindDefaultCredentials(context.Background()); err != nil {
				return
			}
		} else {
			nx := "gs_" + ep + "_url"
			for _, v := range os.Environ() {
				j := strings.Index(v, "=")
				if j > 0 && strings.ToLower(v[:j]) == nx {
					ap, err = defineAccessPoint(v[j+1:])
					if err != nil {
						continue
					}
					break
				}
			}
			if ap == nil {
				err = errors.New("not found access point for gsurl `" + gsurl + "`")
				return
			}
		}
		registry[ep] = ap
	}
	if ap.Bucket != "" {
		path := strings.Trim(ap.Prefix+u.Path, "/")
		loc = Location{ap.Bucket, path}
	} else {
		path := strings.Trim(u.Path, "/")
		j := strings.Index(path, "/")
		if j < 0 {
			err = errors.New("bad gsurl (there is no key name) `" + gsurl + "`")
		}
		loc = Location{path[:j], path[j+1:]}
	}
	return
}

func defineAccessPoint(gsdef string) (ap *AccessPoint, err error) {
	if strings.HasPrefix(gsdef, "json://") {
		s := strings.Split(gsdef[7:], ":")
		if len(s) != 4 {
			return nil, errors.New("invalid access pointdefinition (bad cound of elements)")
		}
		var f *os.File
		var dat []byte
		f, err = os.Open(s[3])
		if err != nil {
			return nil, errors.Wrap(err, "can't open credentials file")
		}
		dat, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read credentials file")
		}
		if s[2] != "" {
			// encrypted
			dat, err = fu.Decrypt(s[2], dat)
			if err != nil {
				return nil, errors.Wrap(err, "failed to decrypt credentials file")
			}
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to read credentials file")
		}
		ap = &AccessPoint{Bucket: s[0], Prefix: s[1]}
		ap.Credentials, err = google.CredentialsFromJSON(context.Background(), dat, storage.DevstorageReadWriteScope)
		return
	}
	return nil, errors.New("invalid access pointdefinition (no json:// prefix)")
}

func (ap *AccessPoint) Service() (svc *storage.Service, err error) {
	mu.Lock()
	defer mu.Unlock()
	if ap.service != nil {
		return ap.service, nil
	}
	httpClient := oauth2.NewClient(context.Background(), ap.Credentials.TokenSource)
	ap.service, err = storage.New(httpClient)
	return ap.service, err
}
