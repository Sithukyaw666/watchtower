package container

import (
	"regexp"
	"sort"

	"github.com/Masterminds/semver/v3"
)

func SelectBestTag(tags []string, strategy string, filter string) (string, error) {
	var filteredTags []string
	if filter != "" {
		re, err := regexp.Compile(filter)
		if err != nil {
			return "", err
		}
		for _, tag := range tags {
			if re.MatchString(tag) {
				filteredTags = append(filteredTags, tag)
			}
		}
	} else {
		filteredTags = tags
	}
	if len(filteredTags) == 0 {
		return "", nil
	}

	switch strategy {
	case "semver":
		return selectHighestSemVer(filteredTags), nil
	// can add new strategy but need to find the logical approach for custom tag
	default:
		return "", nil
	}
}

func selectHighestSemVer(tags []string) string {
	var versions []*semver.Version

	rawTagMap := make(map[string]*semver.Version)

	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		versions = append(versions, v)
		rawTagMap[tag] = v
	}

	if len(versions) == 0 {
		return ""
	}

	sort.Sort(semver.Collection(versions))

	highestVersion := versions[len(versions)-1]

	for tag, v := range rawTagMap {
		if v.Equal(highestVersion) {
			return tag
		}
	}
	return ""
}
