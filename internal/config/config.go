package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Host        string
	Port        string
	StaticDir   string
	TemplateDir string
}

func Load() (*Config, error) {
	envMap, err := parseDotEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:        getEnv(envMap, "HOST", "0.0.0.0"),
		Port:        getEnv(envMap, "PORT", "8080"),
		StaticDir:   getEnv(envMap, "STATIC_DIR", "web/static"),
		TemplateDir: getEnv(envMap, "TEMPLATE_DIR", "web/templates"),
	}, nil
}

func parseDotEnv() (map[string]string, error) {
	env := make(map[string]string)

	file, err := os.Open(".env")
	if err != nil {
		if os.IsNotExist(err) {
			return env, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}

		env[key] = value
	}
	return env, scanner.Err()
}

func getEnv(env map[string]string, key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if value, ok := env[key]; ok {
		return value
	}
	return defaultValue
}
