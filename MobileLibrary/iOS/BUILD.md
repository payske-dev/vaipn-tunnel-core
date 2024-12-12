# Building the Vaipn iOS Library

## Manual Build

### Prerequisites

* xcode `xcode-select --install`

* [git](https://git-scm.com/download/mac)

* [homebrew](http://brew.sh/)

* golang
  - `brew install go`

### Build Steps

* Run `build-vaipn-framework.sh`.
  - If this fails, especially in the `gomobile` step, try re-running it.

* The result will be in `MobileLibrary/iOS/build`.

#### Testing

Run `test-vaipn-framework.sh`.
