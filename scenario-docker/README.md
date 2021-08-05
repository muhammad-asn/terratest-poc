# Docker Scenario Unit Test

1. Check docker port binding, volume binding, container name
2. Check docker version of images
3. Check env hardcoded / missing

## How to use
1. Docker and Go programming language already installed
2. Init and install the dependencies
    ``` bash
    go mod init wordpress-test
    go get github.com/gruntwork-io/terratest 
    ```
3. Execute the test (for example 10m)
    ```
    go test -v -timeout 10m
    ```

Reference:
- https://www.digitalocean.com/community/tutorials/how-to-install-wordpress-with-docker-compose
- https://benmatselby.dev/post/terratest/
- https://stackoverflow.com/questions/61860421/programmatically-exec-into-docker-container
- https://stackoverflow.com/questions/31438112/bash-docker-exec-file-redirection-from-inside-a-container
- https://stackoverflow.com/questions/39478181/how-to-get-output-from-stdout-into-a-string-in-golang