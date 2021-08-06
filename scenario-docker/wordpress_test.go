package test

import (
	"log"
	"net/http"
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/stretchr/testify/assert"
)

func TestUnit(t *testing.T) {

	/**
	Notes:
	For the volume binding test you need to check by the temporary name that built
	by with substring from "<func_name>", for example here "TestUnit" but you should do it
	by lowercase the string

	For example:

	...
	volumes:
		vol_data: /data
	....

	So you need to check the binding volume with <func_name>_vol_data with "testunit_vol_data"
	**/

	dockerComposeFile := "./docker-compose.yml"

	// Run the docker compose file
	docker.RunDockerCompose(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"up",
		"-d",
	)

	dbInstance := docker.Inspect(t, "db")
	wordpressInstance := docker.Inspect(t, "wordpress")

	/**
		MySQL Database Container
	**/
	// Check db port
	assert.Empty(t, dbInstance.Ports)

	// Check db volume binding
	assert.Equal(t, dbInstance.Binds[0].Destination, "testunit_db_data")

	// Check db version
	dbVersion, _ := exec.Command("/bin/sh", "-c", "docker exec -t db /bin/sh -c \"mysql --version\"").Output()
	assert.Contains(t, string(dbVersion), "8.0")

	// Check db environment variable
	mysqlRootPasswordEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t db /bin/sh -c \"env | grep MYSQL_ROOT_PASSWORD\"").Output()
	assert.Contains(t, string(mysqlRootPasswordEnv), "somewordpress")

	mysqlDatabaseEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t db /bin/sh -c \"env | grep MYSQL_DATABASE\"").Output()
	assert.Contains(t, string(mysqlDatabaseEnv), "wordpress")

	mysqlUserEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t db /bin/sh -c \"env | grep MYSQL_USER\"").Output()
	assert.Contains(t, string(mysqlUserEnv), "wordpress")

	mysqlPasswordEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t db /bin/sh -c \"env | grep MYSQL_PASSWORD\"").Output()
	assert.Contains(t, string(mysqlPasswordEnv), "wordpress")

	/**
		Wordpress Container
	**/

	// Check wordpress port and port used
	assert.NotEmpty(t, wordpressInstance.Ports)
	assert.Equal(t, wordpressInstance.Ports[0].HostPort, uint16(8000))

	// Check wordpress volume binding
	assert.Equal(t, wordpressInstance.Binds[0].Destination, "testunit_wordpress_data")

	// Check Wordpress version
	wordpressVersion, _ := exec.Command("/bin/sh", "-c", "docker exec -t wordpress /bin/sh -c \"grep wp_version wp-includes/version.php\"").Output()
	assert.Contains(t, string(wordpressVersion), "5.8")

	// Check wordpress environment variable
	wordpressDbHostEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t wordpress /bin/sh -c \"env | grep WORDPRESS_DB_HOST\"").Output()
	assert.Contains(t, string(wordpressDbHostEnv), "db:3306")

	wordpressDbUserEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t wordpress /bin/sh -c \"env | grep WORDPRESS_DB_USER\"").Output()
	assert.Contains(t, string(wordpressDbUserEnv), "wordpress")

	wordpressDbPasswordEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t wordpress /bin/sh -c \"env | grep WORDPRESS_DB_PASSWORD\"").Output()
	assert.Contains(t, string(wordpressDbPasswordEnv), "wordpress")

	wordpressDbNameEnv, _ := exec.Command("/bin/sh", "-c", "docker exec -t wordpress /bin/sh -c \"env | grep WORDPRESS_DB_NAME\"").Output()
	assert.Contains(t, string(wordpressDbNameEnv), "wordpress")

	// Cleaning
	defer docker.RunDockerComposeAndGetStdOut(t, &docker.Options{}, "-f", dockerComposeFile, "down")

}

func TestFunctionalWordpress(t *testing.T) {
	// TODO: Check functionality of Wordpress service

	dockerComposeFile := "./docker-compose.yml"
	docker.RunDockerCompose(
		t,
		&docker.Options{},
		"-f",
		dockerComposeFile,
		"up",
		"-d",
	)

	resp, err := http.Get("http://localhost:8000/")

	if err != nil {
		log.Fatal(err)
	}

	installUrl := resp.Request.URL.String()

	// Redirect and check if go to install page, (200)
	assert.Equal(t, installUrl, "http://localhost:8000/wp-admin/install.php")
	assert.Equal(t, resp.StatusCode, 200)

	defer resp.Body.Close()
	defer docker.RunDockerComposeAndGetStdOut(t, &docker.Options{}, "-f", dockerComposeFile, "down")

}
