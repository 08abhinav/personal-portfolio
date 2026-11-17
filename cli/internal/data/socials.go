package data

import "strings"

type Social struct {
	Name string
	URL  string
}

var Socials = []Social{
	{"github", "https://github.com/08abhinav"},
	{"linkedin", "https://linkedin.com/in/08abhinav"},
	{"medium", "https://medium.com/@abhinavnegi101"},
	{"email", "mailto:abhinavnegi101@gmail.com"},
}

func FindSocial(name string) (Social, bool) {
	for _, s := range Socials {
		if strings.EqualFold(s.Name, name) {
			return s, true
		}
	}
	return Social{}, false
}

func SocialNames() []string {
	names := make([]string, len(Socials))
	for i, s := range Socials {
		names[i] = s.Name
	}
	return names
}
