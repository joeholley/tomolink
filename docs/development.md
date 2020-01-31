# Tomolink Development Guide
  
This guide is for developers wanting to implement changes to the Tomolink source code. This guide only covers how to make changes in conjunction with Google Cloud Platform products and the GCP development experience. 
  
**If you are a user of Tomolink and don't plan to make changes to the source code or recompile, then you likely don't need this guide. Please refer to the [README](../README.md) file for more information on using Tomolink**

## Local development

Tomolink can be run locally in your development environment, using the Firestore emulator provided in the Google Cloud `gcloud` command line suite ([installation instructions](https://cloud.google.com/sdk/install)).
1. Clone the source repository, if you haven't already.  It is recommended you clone it to your local golang `src/` directory, but you can also set your [GOPATH](https://github.com/golang/go/wiki/GOPATH) instead if you feel comfortable doing so.
```bash
#TODO: update with final github URI
git clone https://github.com/joeholley/tomolink.git
```
It is recommended that you use the provided development config file when running Tomolink locally:
```bash
~/go/src/github.com/joeholley/tomolink$ ln -sf ../internal/config/local_dev.yaml cmd/tomolink_defaults.yaml
```
2. Run a copy of the Firestore Emulator on a port that won't conflict with the default Tomolink port of 8080.  This example uses 8081.
```bash
gcloud beta emulators firestore start --project localemu --host-port "localhost:8081"
```
3. Run Tomolink in a separate terminal. By default the Firestore client libraries Tomolink uses will check the `FIRESTORE_EMULATOR_HOST` environment variable.
```bash
cd cmd
export FIRESTORE_EMULATOR_HOST=localhost:8081 DATABASE_ID=localemu; go run httpserver.go 
```
You should see output like the following (your output will depend on your config settings):
```
INFO[0000] Environment variables will override default YAML goconfig values if both exist. Due to differences between how environment variables and YAML goconfig keys are named, please 
see the documentation for more details on this feature  application=tomolink component=internal.config
INFO[0000] Override for default YAML value found in environment variable!  application=tomolink component=internal.config envVarName=DATABASE_ID overrideValue=localemu yamlConfigKey=databas
e.id
INFO[0000] Read config defaults from ./tomolink_defaults.yaml with overrides from env vars  application=tomolink component=internal.config
INFO[0000] No 'relationships.definitions.4.name' configured. Done processing relationships from the config  application=tomolink component=internal.config
INFO[0000] [Dev] Logging all goconfiguration settings    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.1.type = score  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.9 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.3.name = blocks  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  database.engine = firestore            application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  http.port = 8080                       application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.1.name = influencers  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.0.type = score  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  database.id = test                     application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  dev = true                             application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  logging.level = debug                  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  database.options.grpc.pool = 20        application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.strict = false           application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.2.type = score  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  logging.verbose = true                 application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.5 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  http.gracefulwait = 15                 application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  logging.format = text                  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  http.request.readLimit = 500           application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.6 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.4 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.0.name = friends  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.7 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.2.name = followers  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.3.type = score  application=tomolink component=internal.config dev=true
INFO[0000] [Dev]  relationships.definitions.8 = <nil>    application=tomolink component=internal.config dev=true
INFO[0000] database option set                           application=tomolink component=internal.config database.engine=firestore dev=true grpc.pool=20
WARN[0000] Trace logging level configured. Not recommended for production! 
WARN[0000] Verbose logging configured. Not recommended for production! 
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:99 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route="/users/{UUIDSource}/{relationship}/{UUIDTarget}"
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:109 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route="/users/{UUIDSource}/{relationship}"
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:123 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route="/users/users/{UUIDSource}"
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:125 github.com/joeholley/tomolink/internal/app/tomolink.Router() All configured relationship endpoints created 
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:136 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route=/createRelationship
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:147 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route=/updateRelationship
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/internal/app/tomolink/router.go:158 github.com/joeholley/tomolink/internal/app/tomolink.Router() Added route                                   name=todo route=/deleteRelationship
INFO[0000]/usr/local/google/home/joeholley/go/src/github.com/joeholley/tomolink/cmd/httpserver.go:101 main.main() Starting HTTP server                          app=tomolink component=app.main port=8080

```
If all is working, you should see the following in the Firestore emulator terminal window when Tomolink connects:
```
[firestore] INFO: Detected HTTP/2 connection.
```
4. You can now use `curl` or your browser to access Tomolink at [http://127.0.0.1:8080](http://127.0.0.1:8080).

**Note: Tomolink doesn't serve anything at the base URI (http://127.0.0.1:8080).  You'll need to specify the URI that corresponds to an API call to successfully access Tomolink.  For more details, see the [User Guide](userguide.md).**
  
## Building Tomolink
See the [build and deploy guide](builddeploy.md).
