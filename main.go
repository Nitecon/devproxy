package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v2"
)

var config Config

func setLogger() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		return
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

type Server struct {
	Name string   `yaml:"name"`
	Path []string `yaml:"path"`
	Port int      `yaml:"port"`
}

type Config struct {
	ListenPort  int      `yaml:"listen_port"`
	DefaultPort int      `yaml:"default_port"`
	Servers     []Server `yaml:"servers"`
}

func RouterHandler(w http.ResponseWriter, r *http.Request) {
	// loop through the servers and match the path, and add the default server that contains / path to the end of the list
	for _, server := range config.Servers {
		for _, path := range server.Path {
			plen := len(path)
			// if r.Url.Path starts with path, then redirect to the server
			if r.URL.Path[:plen] == path {
				body, headers, status := SendRequest(server, r)
				// set the headers
				for key, value := range headers {
					for _, v := range value {
						w.Header().Add(key, v)
					}
				}
				// write the response body
				w.Write(body)
				log.Info().Msgf("[%s] - (%s:%d) - %s", server.Name, r.Method, status, r.RequestURI)
				return
			}
		}
	}
	// ifno server is found SendRequest to the deafault server
	default_server := Server{Port: config.DefaultPort, Name: "default"}
	body, headers, status := SendRequest(default_server, r)
	// set the headers
	for key, value := range headers {
		for _, v := range value {
			w.Header().Add(key, v)
		}
	}
	// write the response body
	w.Write(body)
	log.Info().Msgf("[%s] - (%s:%d) - %s", "default", r.Method, status, r.RequestURI)
}

func SendRequest(server Server, r *http.Request) (body []byte, headerMap map[string][]string, status int) {
	// create a new http client
	client := &http.Client{}
	// create a new request
	req, err := http.NewRequest(r.Method, fmt.Sprintf("http://localhost:%d%s", server.Port, r.URL.Path), r.Body)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to create a new request to %s", server.Name)
	}
	// copy the headers
	for key, value := range r.Header {
		req.Header.Set(key, value[0])
	}
	// send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send the request to %s", server.Name)
	}
	// read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to read the response body from %s", server.Name)
	}

	headerMap = make(map[string][]string)
	// get the response headers
	for key, value := range resp.Header {
		headerMap[key] = value
	}
	// return the response body and status code
	return body, headerMap, resp.StatusCode
}

func main() {
	setLogger()
	// Read the YAML file into a Config struct.
	yamlData, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read YAML file")
	}

	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal YAML data")

	}

	// Create a ServeMux to handle incoming requests.
	mux := http.NewServeMux()
	mux.HandleFunc("/", RouterHandler)

	// Start the server.
	port := fmt.Sprintf(":%d", config.ListenPort)
	log.Info().Msgf("Listening on port %s...", port)
	log.Fatal().Err(http.ListenAndServe(port, mux))
}
