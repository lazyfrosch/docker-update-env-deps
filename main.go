package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
	"update-deps/dockerfile"
	"update-deps/repos"
)

type updateInfo struct {
	Key     string
	Current string
	Latest  string
}

// TODO: config
var repositoryMap = map[string]string{
	"ICINGAWEB_VERSION":             "https://github.com/Icinga/icingaweb2.git",
	"ICINGA_PHP_LIBRARY_VERSION":    "https://github.com/Icinga/icinga-php-library.git",
	"ICINGA_PHP_THIRDPARTY_VERSION": "https://github.com/Icinga/icinga-php-thirdparty.git",
	"ICINGA_ICINGADB_VERSION":       "https://github.com/Icinga/icingadb-web.git",
	"ICINGA_DIRECTOR_VERSION":       "https://github.com/Icinga/icingaweb2-module-director.git",
	"ICINGA_FILESHIPPER_VERSION":    "https://github.com/Icinga/icingaweb2-module-fileshipper.git",
	"ICINGA_IPL_VERSION":            "https://github.com/Icinga/icingaweb2-module-ipl.git",
	"ICINGA_INCUBATOR_VERSION":      "https://github.com/Icinga/icingaweb2-module-incubator.git",
	"ICINGA_REACTBUNDLE_VERSION":    "https://github.com/Icinga/icingaweb2-module-reactbundle.git",
}

var reEnvVersion = regexp.MustCompile(`(?i)version`)

func main() {
	environmentVars := dockerfile.LookupEnvironment("Dockerfile")

	updates, err := checkUpdates(environmentVars, repositoryMap)
	if err != nil {
		log.Fatal(err)
	}

	for _, update := range updates {
		log.WithFields(log.Fields{
			"key":     update.Key,
			"current": update.Current,
			"latest":  update.Latest,
		}).Warnf("found update for %s: %s -> %s", update.Key, update.Current, update.Latest)
	}

	if len(updates) > 0 {
		os.Exit(1)
	}
}

func checkUpdates(current map[string]string, repos map[string]string) (updates []updateInfo, err error) {
	for key, version := range current {
		url, found := repos[key]

		if !found {
			if reEnvVersion.MatchString(key) {
				log.WithField("key", key).Warn("variable looks like a version, but is not mapped")
			}

			continue
		}

		logInfo := log.WithFields(log.Fields{
			"key":     key,
			"version": version,
		})

		logInfo.Debug("found current version")

		var latest string
		latest, err = latestGitHub(url)
		if err != nil {
			return
		}

		if version != latest {
			info := updateInfo{
				Key:     key,
				Current: version,
				Latest:  latest,
			}
			updates = append(updates, info)
		} else {
			logInfo.Debug("version is up to date")
		}
	}

	return
}

func latestGitHub(url string) (latest string, err error) {
	repo, err := repos.LoadGitHub(url)
	if err != nil {
		err = fmt.Errorf("could not load repo: %w", err)
		return
	}

	releases, err := repo.LoadReleases()
	if err != nil {
		err = fmt.Errorf("could not load releases: %w", err)
		return
	}

	latest = normalizeVersion(releases[0])
	return
}

func normalizeVersion(version string) string {
	version = strings.TrimPrefix(version, "v")

	return version
}
