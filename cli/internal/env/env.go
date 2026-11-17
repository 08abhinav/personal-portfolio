package env

import (
	"time"

	"github.com/08abhinav/personal-portfolio/cli/internal/geo"
)

const (
	ModeLocal 	= "local"
	ModeWeb 	= "web"
)

type Client struct{
	IP			string
	UserAgent 	string
	Timezone 	string
}
type Env struct{
	Mode 		string
	StartedAt 	time.Time
	Client		Client
	Geo			*geo.Resolver
}

func Local() *Env{
	return &Env{Mode: ModeLocal, StartedAt: time.Now()}
}
