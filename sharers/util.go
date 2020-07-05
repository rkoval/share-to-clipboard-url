package sharers

import "regexp"

func FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		if i != 0 && name != "" {
			results[regex.SubexpNames()[i]] = name
		}
	}
	return results
}
