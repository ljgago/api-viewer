# API Viewer

An OpenAPI Specification (OAS) viewer using [Scalar](https://github.com/scalar/scalar).

## Features

- Text and JSON logs with request and response headers.
- Internal proxy server to bypass CORS.
- Scalar UI with theme setting.

## Usage

```sh
$ api-viewer -h
An OpenAPI Specification (OAS) viewer

Usage:
  api-viewer [flags]

Flags:
  -f, --file string    spec file
  -h, --help           help for api-viewer
  -j, --json           show log format in json
  -p, --port string    server port (default "8080")
  -x, --proxy string   proxy port (default "9090")
  -t, --theme string   color theme (default "default")
  -v, --version        version of api-viewer
```

## Themes

The list of Scalar themes are:

```js
[
  "alternate" | "default" | "moon" | "purple" |
  "solarized" | "bluePlanet" | "saturn" | "kepler" |
  "mars" | "deepSpace" | "laserwave" | "none"
]
```

## License
[MIT License](./LICENSE)
