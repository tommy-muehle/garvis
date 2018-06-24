package configuration

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	languages = []string{"PHP", "GO"}
)

type Configuration struct {
	Server    serverConfiguration     `yaml:"server"`
	Log       logConfiguration        `yaml:"log"`
	Github    githubConfiguration     `yaml:"github"`
	Reviewers []reviewerConfiguration `yaml:"reviewers"`
}

type serverConfiguration struct {
	Addr string `yaml:"addr"`
}

type logConfiguration struct {
	Debug bool `yaml:"debug"`
}

type githubConfiguration struct {
	SecretKey string `yaml:"secrect_key"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type reviewerConfiguration struct {
	Language string   `yaml:"language"`
	Users    []string `yaml:"users"`
}

func LoadFromFile(file string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Configuration{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, err
	}

	return config, nil
}

func MustLoadFromFile(file string) *Configuration {
	config, err := LoadFromFile(file)
	if err != nil {
		panic(err)
	}

	return config
}

func DefaultConfig() *Configuration {
	return &Configuration{
		Server: serverConfiguration{
			Addr: ":80",
		},
		Log: logConfiguration{
			Debug: true,
		},
		Github: githubConfiguration{
			SecretKey: "secret",
			Username:  "bot-user",
			Password:  "password",
		},
		Reviewers: []reviewerConfiguration{
			reviewerConfiguration{
				Language: "PHP",
				Users: []string{
					"username1",
					"username2",
				},
			},
			reviewerConfiguration{
				Language: "Go",
				Users: []string{
					"username3",
					"username4",
				},
			},
		},
	}
}

func FromEnvironment() *Configuration {
	debug := true
	b, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		debug = b
	}

	c := DefaultConfig()
	c.Server.Addr = ":" + os.Getenv("PORT")
	c.Log.Debug = debug

	c.Github.Password = os.Getenv("GITHUB_PASSWORD")
	c.Github.SecretKey = os.Getenv("GITHUB_SECRET")
	c.Github.Username = os.Getenv("GITHUB_USERNAME")

	for _, language := range languages {
		reviewers := os.Getenv(language + "_REVIEWERS")
		c.Reviewers = append(c.Reviewers, reviewerConfiguration{
			Language: language,
			Users:    strings.Split(reviewers, ","),
		})
	}

	return c
}
