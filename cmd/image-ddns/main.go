package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/caarlos0/env"
	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/gofrs/uuid"
	"github.com/otiai10/gosseract"
	"github.com/sirupsen/logrus"
)

type config struct {
	CloudflareApiKey   string `env:"CLOUDFLARE_API_KEY"`
	CloudflareApiEmail string `env:"CLOUDFLARE_API_EMAIL"`
	Name               string `env:"NAME"`
	Zone               string `env:"ZONE"`
	ImageUrl           string `env:"IMAGE_URL"`
	LogLevel           string `env:"LOG_LEVEL" envDefault:"info"`
	TmpDir             string `env:"TMP_DIR" envDefault:"/tmp"`
}

const (
	rtype = "A"
	ttl   = 1
	proxy = true
)

var (
	cfg          config
	ipPortRegexp *regexp.Regexp
)

func downloadUrl(uri string) (string, error) {
	var err error

	newUUID := uuid.Must(uuid.NewV4())

	tmpFilename := filepath.Join(cfg.TmpDir, newUUID.String())

	out, err := os.Create(tmpFilename)
	if err != nil {
		return "", err
	}

	defer out.Close()

	logrus.Debugf("downloading %s to %s", uri, tmpFilename)

	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	io.Copy(out, resp.Body)

	return tmpFilename, nil
}

func checkError(label string, err error) {
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf(label)
		panic(err)
	}
}

func updateCloudflare(externalIP string) {
	// Construct a new API object
	api, err := cloudflare.New(cfg.CloudflareApiKey, cfg.CloudflareApiEmail)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("new cloudflare")
		panic(err)
	}

	zoneID, err := api.ZoneIDByName(cfg.Zone)
	if err != nil {
		logrus.Errorf("Error updating DNS record: %#+v\n", err)
		return
	}

	// Look for an existing record
	rr := cloudflare.DNSRecord{
		Name: cfg.Name + "." + cfg.Zone,
	}
	records, err := api.DNSRecords(zoneID, rr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("fetching DNS records")
		return
	}

	if len(records) > 0 {
		for _, r := range records {
			if r.Type == rtype {
				rr.ID = r.ID
				rr.Type = r.Type
				rr.Content = externalIP
				rr.TTL = ttl
				rr.Proxied = proxy
				err := api.UpdateDNSRecord(zoneID, r.ID, rr)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"err": err,
					}).Errorf("updating DNS record")
				} else {
					logrus.Debugf("Updated %s.%s -> %s\n", cfg.Name, cfg.Zone, externalIP)
				}
			}
		}
	} else {
		rr.Type = rtype
		rr.Content = externalIP
		rr.TTL = ttl
		rr.Proxied = proxy
		_, err = api.CreateDNSRecord(zoneID, rr)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("creating DNS record")
		} else {
			logrus.Debugf("Added %s.%s -> %s\n", cfg.Name, cfg.Zone, externalIP)
		}
	}
}

func main() {
	var err error

	cfg = config{}
	err = env.Parse(&cfg)
	checkError("config parse", err)

	level, err := logrus.ParseLevel(cfg.LogLevel)
	checkError("parse log level", err)
	logrus.SetLevel(level)

	imageFileName, err := downloadUrl(cfg.ImageUrl)
	checkError("downloadUrl", err)
	defer func() {
		logrus.Debugf("removing file %s", imageFileName)
		os.Remove(imageFileName)
	}()

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(imageFileName)
	ocrText, err := client.Text()
	checkError("ocr text", err)

	logrus.WithFields(logrus.Fields{"ocrText": ocrText}).Debugf("ocr text")

	matches := ipPortRegexp.FindStringSubmatch(ocrText)
	parts := strings.Split(strings.ReplaceAll(matches[0], ",", "."), ":")
	logrus.WithFields(logrus.Fields{"matches": matches, "parts": parts}).Debugf("find string submatch")
	updateCloudflare(parts[0])
}

func init() {
	ipPortRegexp = regexp.MustCompile("([0-9,.]+:[0-9]+)")
}
