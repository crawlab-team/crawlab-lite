package services

import (
	"crawlab-lite/results"
	"github.com/apex/log"
	"github.com/imroc/req"
	"runtime/debug"
	"sort"
)

func GetLatestRelease() (release results.Release, err error) {
	res, err := req.Get("https://api.github.com/repos/crawlab-team/crawlab-lite/releases")
	if err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return release, err
	}

	var releaseDataList results.ReleaseSlices
	if err := res.ToJSON(&releaseDataList); err != nil {
		log.Errorf(err.Error())
		debug.PrintStack()
		return release, err
	}

	if len(releaseDataList) == 0 {
		return results.Release{}, err
	}

	sort.Sort(releaseDataList)

	return releaseDataList[len(releaseDataList)-1], nil
}
