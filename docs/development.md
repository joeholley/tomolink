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
  
Tomolink is a Golang project and uses the [native module support in Golang >1.13 as it's dependancy management system](https://blog.golang.org/using-go-modules). Tomolink supports a Docker container build target, and is built using GCP's [Cloud Build](https://cloud.google.com/cloud-build/) service to build it. The necessary [Dockerfile](../Dockerfile) and [cloudbuild.yaml](../cloudbuild.yaml) files are contained in the root of the repository.  

Building Tomolink using GCP uses the following paid products:
* [Container Registry](https://cloud.google.com/container-registry)
* [Cloud Build](https://cloud.google.com/cloud-build/docs)

Deploying Tomolink uses the following paid products:
* [Cloud Run](https://cloud.google.com/run/docs)
* [Firestore](https://cloud.google.com/firestore/docs)

These products cab incur charges, although typical usage would be eligible for the [free tier](https://cloud.google.com/free).

### Before you begin
1. Clone the source repository, if you haven't already.
```bash
#TODO: update with final github URI
git clone https://github.com/joeholley/tomolink.git
```
2. If you don't have a `gcloud` command line environment set up and Cloud Build enabled on a Google Cloud Project, complete the ["Before You Begin"](https://cloud.google.com/cloud-build/docs/quickstart-docker#before-you-begin) and ["Log in to Google Cloud"](https://cloud.google.com/cloud-build/docs/quickstart-docker#log_in_to) steps from the [Cloud Build Quickstart for Docker](https://cloud.google.com/cloud-build/docs/quickstart-docker#log_in_to).  Stop when you get to the "Preparing source files" section of the guide, as you completed that in step #1.
3. If you haven't already, enable the GCP [Container Registry API](https://console.cloud.google.com/flows/enableapi?apiid=containerregistry.googleapis.com&redirect=https://cloud.google.com/container-registry/docs/quickstart&_ga=2.79131540.1852199750.1580108081-50451190.1521501879) .


### Using Cloud Build to build and deploy Tomolink

**Note: each project only supports one Firestore instance. Although Tomolink can be run in a project with another application using Firestore as long as they don't use the same documents, it is not recommended.**
1. Set up a Firestore instance in your project by following the [Create a Cloud Firestore in Native mode Database](https://cloud.google.com/firestore/docs/quickstart-servers#create_a_in_native_mode_database) part of the quickstart guide.   As long as you will deploy Tomolink to the same GCP project as the Firestore instance, Tomolink will automatically be able to access Firestore to save data without additional authentication.
2. If you haven't already, enable the GCP [Cloud Run API](http://console.cloud.google.com/apis/library/run.googleapis.com?_ga=2.23508026.1852199750.1580108081-50451190.1521501879)
3. Start the build by running the following command from the repository root. It both builds the container and deploys it to Google Cloud Run. It will take a moment to complete.
```bash
gcloud builds submit --config cloudbuild.yaml .
```
4. You can view your deployment by going to [Cloud Run in the Google Cloud Console](https://console.cloud.google.com/run?enableapi=true&_ga=2.224795418.1852199750.1580108081-50451190.1521501879).
5. Finally, if you want to allow access to Tomolink from outside your GCP project, you'll need to [allow unauthenticated access](https://cloud.google.com/run/docs/authenticating/public).  **This is not recommended for production.**

### Setting up Continuous Deployment of Tomolink
You can use Cloud Build to automate builds and deployments to Cloud Run. See the guide on [Continuous Deployment from git using Cloud Build](https://cloud.google.com/run/docs/continuous-deployment-with-cloud-build).
