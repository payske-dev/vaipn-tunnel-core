## Vaipn Android Library README

### Overview

Vaipn Library for Android enables you to easily embed Vaipn in your Android
app.

### Using the Library

#### If building from source

 1. Build `ca.vaipn.aar` from via the docker container.
 2. Add `ca.vaipn.aar` to your Android Studio project as described in the [gomobile documentation](https://godoc.org/golang.org/x/mobile/cmd/gomobile)

#### If using Maven based binary distribution

1. Add maven repo to your app build.gradle
```
repositories {
    ...
    maven {
        url "https://raw.github.com/payske-dev/vaipn-tunnel-core-Android-library/master"
    }
}
```
then add VaipnTunnel dependency like following
```
dependencies {
    ...
    implementation 'ca.vaipn:vaipntunnel:x.y.z'
}
```
Where x.y.z is the target version. Latest available release version can be found at https://github.com/payske-dev/vaipn-tunnel-core-Android-library

See example usage in [TunneledWebView sample app](./SampleApps/TunneledWebView/README.md)

### Building with Docker

Note that you may need to use `sudo docker` below, depending on your OS.

##### Create the build image:

1. While in the `MobileLibrary/Android` directory, run the command: `docker build --no-cache=true -t psiandroid -f Dockerfile ..`

2. Once completed, verify that you see an image named `psiandroid` when running: `docker images`

##### Run the build:

*Ensure that the command below is run from within the `MobileLibrary/Android` directory*

*The `--user "$(id -u):$(id -g)"` flag must be omitted when building on macOS because manually specifying the UID/GID will cause the build to fail since Docker on macOS handles host to container UID/GID mapping automatically and will be unable to map the manually specified UID/GID to the current user.*

```bash
cd ../.. && \
  docker run \
  --platform=linux/amd64 \
  --user "$(id -u):$(id -g)" \
  --rm \
  -v $(pwd):/go/src/github.com/payske-dev/vaipn-tunnel-core \
  psiandroid \
  /bin/bash -c 'cd /go/src/github.com/payske-dev/vaipn-tunnel-core/MobileLibrary/Android && ./make.bash' \
; cd -
```

When that command completes, the compiled `.aar` files (suitable for use in an Android Studio project) will be located in the current directory (it will likely be owned by root, so be sure to `chown` to an appropriate user).

### Building without Docker (from source)

##### Prerequisites:

 - The `build-essential` package (on Debian based systems - or its equivalent for your platform)
 - Go 1.15 or later
 - Full JDK
 - Android NDK
 - Android SDK

##### Steps:

 1. Follow Go Android documentation ([gomobile documentation](https://godoc.org/golang.org/x/mobile/cmd/gomobile))
 2. Run `make.bash`
