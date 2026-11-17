package data


type Profile struct {
	Name   string   `json:"name"`
	Title  string   `json:"title"`
	Role   string   `json:"role"`
	Focus  []string `json:"focus"`
	Status string   `json:"status"`
	Site   string   `json:"site"`
	Next   string   `json:"next"`
}

var Me = Profile{
	Name:   "Abhinav Negi",
	Title:  "Backend & cloud-native developer",
	Role:   "Postgrad in computer application",
	Focus:  []string{"Go", "Kubernetes", "AWS"},
	Status: "Open to backend, cloud and DevOps roles",
	Site:   "https://abhinavnegi.site",
	Next:   "socials",
}
