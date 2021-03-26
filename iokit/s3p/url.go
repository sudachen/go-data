package s3p

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"net/url"
	"strings"
	"sudachen.xyz/pkg/go-data/errors"
)

/*
	S3_*_URL => ...
	S3_URL = s3://key1:secret1@region1.entrypoint/prefix1
	s3://$/xxx => Lookup("") => AccessPoint{entrypoint,region1,prefix1,{key1,secret1}} + xxx
	S3_DEFAULT_URL = s3://key2:secret2@region2.entrypoint/prefix2
	s3://$default/xxx => Lookup("default") => AccessPoint{entrypoint,region2,prefix2,{key2,secret2}} + xxx
*/

func DecodeUrl(u *url.URL) (ap *AccessPoint, err error) {
	p := u.Path
	for len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	j := strings.Index(p, "/")
	ap = &AccessPoint{}
	if j < 0 {
		ap.Bucket = p
	} else {
		ap.Bucket = p[:j]
		ap.Prefix = p[j+1:]
	}
	if ap.Bucket == "" {
		return nil, errors.New("bad bucket name in path `" + u.Path + "`")
	}
	hs := strings.Split(u.Host, ".")
	if len(hs) > 2 {
		if hs[0] == "s3" {
			hs = hs[1:]
		}
		ap.Region = hs[0]
	}
	ap.Endpoint = u.Host
	if pwd, ok := u.User.Password(); !ok {
		ap.Credentials = credentials.NewCredentials(&credentials.SharedCredentialsProvider{})
	} else {
		ap.Credentials = credentials.NewStaticCredentials(u.User.Username(), pwd, "")
	}
	return
}

func ResolveUrl(s3url string) (ssn *session.Session, loc Location, err error) {
	u, err := url.Parse(s3url)
	if err != nil {
		return
	}
	if len(u.Host) > 0 && u.Host[0] == '$' {
		apname := u.Host[1:]
		ap, ok := Lookup(apname)
		if !ok {
			return nil, Location{}, errors.New("unknown access point " + apname)
		}
		ssn, err = ap.Session()
		path := strings.Trim(ap.Prefix+u.Path, "/")
		loc = Location{aws.String(ap.Bucket), aws.String(path)}
		return
	}
	ap, err := DecodeUrl(u)
	if err != nil {
		return
	}
	loc = Location{aws.String(ap.Bucket), aws.String(ap.Prefix)}
	ssn, err = ap.Session()
	return
}

func Download(url string, wr io.WriterAt) (err error) {
	ssn, loc, err := ResolveUrl(url)
	if err != nil {
		return
	}
	dlr := s3manager.NewDownloader(ssn)
	_, err = dlr.Download(wr, &s3.GetObjectInput{Bucket: loc.Bucket, Key: loc.Key})
	return
}

func Upload(url string, rd io.Reader, metadata map[string]string) (err error) {
	ssn, loc, err := ResolveUrl(url)
	if err != nil {
		return
	}
	uploader := s3manager.NewUploader(ssn)
	mdp := map[string]*string{}
	for k, v := range metadata {
		mdp[k] = aws.String(v)
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:   loc.Bucket,
		Key:      loc.Key,
		Body:     rd,
		Metadata: mdp,
	})
	return
}
