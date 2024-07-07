# **Yamlink**

Yamlink is a simple URL mapping and redirecting system that can be used to map URLs in a YAML configuration file. The tool can be used to create custom URLs for internal or external use.

**Features**

- Supports mapping URLs with depth up to 3 (e.g., `key1/0`, `key1/1/key2`)
- Can handle arrays of values in the YAML configuration file
- Automatically refreshes the map when the underlying YAML file changes

## **Example Usage**

In this example, we have a YAML configuration file that maps URLs to their corresponding values:

```yaml
key1:
  - example.com
```

Using the `yamlink.ParseYaml` function, we can parse this YAML file and use it to look up URL paths. For instance, if we request the path `/key1/0`, Yamlink will return `example.com`.

## **Building and Running**

To build and run Yamlink, simply execute the following commands:

```bash
go build server.go
./yamlink
```

This will start the Yamlink server, which can be accessed at `http://localhost:8080`. The server supports two endpoints:

- `/`: Returns a health check response
- `<path>`: Redirects to the corresponding URL value in the YAML configuration file

## **Configuration**

The server requires 2 configuration files to run:

- `config.yaml`: Contains the server configuration
  This file should be in the following format:
  ```yaml
  host: localhost
  port: 8080
  map_file: map.yaml
  refresh_interval: 5 # in seconds
  ```
- `map.yaml`: Contains the URL mapping
  This file should be in the following format:
  ```yaml
  key1:
    - https://example.com
  example: https://example.com
  ```

