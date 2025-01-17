# Jenkins + Docker in Docker + Shipshape

The Dockerfile in this directory can be used to build a Docker image which has
both Docker in Docker and Jenkins. The docker file is based off of the
Dockerfile `shipshape/test/dind/Dockerfile`.

## Prerequisites

You must first build shipshape. Follow the instructions in the top-level README.
This will produce a shipshape binary in `bazel-bin/shipshape/cli/` needed for
this docker image.

## Building

From this directory run this command to build the Jenkins image:

    $ ./build-image.sh -t user/jenkins

This builds the image with the tag `user/jenkins`. You could pick any other tag,
but this demo tag is used in the below examples.

## Running

After building the image you can launch it using for example the command

    $ docker run --privileged -p 8080:8080 user/jenkins

This launches Jenkins and binds its web port to `localhost:8080`.  If you need
to access the Jenkins instance via the [Jenkins CLI][1], you will need to bind
port 50000 as well:

    $ docker run --privileged -p 8080:8080 -p 50000:50000 user/jenkins

If you need persistent storage for your Jenkins workspace, add `-v
/my/jenkins:/var/jenkins_home` to the command:

    $ docker run --privileged -p 8080:8080 -p 50000:50000 \
        -v /home/myjenkins:/var/jenkins_home user/jenkins

this mounts the `/home/myjenkins` directory as `/var/jenkins_home` inside the
Jenkins instance.

## Install Plugin

You can install a plugin with the [Jenkins CLI][1]:

    $ curl http://localhost:8080/jnlpJars/jenkins-cli.jar -o jenkins-cli.jar \
        && java -jar jenkins-cli.jar -s http://localhost:8080/ install-plugin \
        ../../jenkins_plugin/target/shipshape-plugin.hpi -restart

## Local Docker Registry

Running a [local docker registry][2] to host private analyzer images can be done
using these commands:

    $ docker run -d -p 5000:5000 --restart=always --name registry registry:2
    $ docker tag user/analyzer localhost:5000/user/analyzer
    $ docker push localhost:5000/user/analyzer
    $ docker run --privileged --net=host user/jenkins

The new option `--net=host` tells Docker to use the host network stack. This way
your local docker registry can be used in the docker image URI for your
analyzer: `localhost:5000/user/analyzer:latest`. When using `--net=host` it
is redundant to specify port bindings since the container will have access to
all the host ports.

[1]: https://wiki.jenkins-ci.org/display/JENKINS/Jenkins+CLI
[2]: https://docs.docker.com/registry/deploying/
