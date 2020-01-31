# Building Tomolink
  
Tomolink is a Golang project and uses the [native module support in Golang >1.13 as it's dependancy management system](https://blog.golang.org/using-go-modules). Tomolink supports a Docker container build target, and is built using GCP's [Cloud Build](https://cloud.google.com/cloud-build/) service to build it. The necessary [Dockerfile](../Dockerfile) and [cloudbuild.yaml](../cloudbuild.yaml) files are contained in the root of the repository.  

The included [cloudbuild.yaml](../cloudbuild.yaml) file specifies three steps that:
1) build the docker container
2) push the docker container to your project's Google Container Registry
3) deploy your docker container as a new revision on Cloud Run

Building Tomolink using GCP uses the following paid products:
* [Container Registry](https://cloud.google.com/container-registry)
* [Cloud Build](https://cloud.google.com/cloud-build/docs)

Deploying Tomolink uses the following paid products:
* [Cloud Run](https://cloud.google.com/run/docs)
* [Firestore](https://cloud.google.com/firestore/docs)

These products cab incur charges, although typical usage would be eligible for the [free tier](https://cloud.google.com/free).

## Before you begin
1. Clone the source repository, if you haven't already.
```bash
#TODO: update with final github URI
git clone https://github.com/joeholley/tomolink.git
```
2. If you don't have a `gcloud` command line environment set up and Cloud Build enabled on a Google Cloud Project, complete the ["Before You Begin"](https://cloud.google.com/cloud-build/docs/quickstart-docker#before-you-begin) and ["Log in to Google Cloud"](https://cloud.google.com/cloud-build/docs/quickstart-docker#log_in_to) steps from the [Cloud Build Quickstart for Docker](https://cloud.google.com/cloud-build/docs/quickstart-docker#log_in_to).  Stop when you get to the "Preparing source files" section of the guide, as you completed that in step #1.
3. If you haven't already, enable the GCP [Container Registry API](https://console.cloud.google.com/flows/enableapi?apiid=containerregistry.googleapis.com&redirect=https://cloud.google.com/container-registry/docs/quickstart&_ga=2.79131540.1852199750.1580108081-50451190.1521501879) .

## Using Cloud Build to build and deploy Tomolink

**Note: each project only supports one Firestore instance. Although Tomolink can be run in a project with another application using Firestore as long as they don't use the same documents, it is not recommended.**
1. Set up a Firestore instance in your project by following the [Create a Cloud Firestore in Native mode Database](https://cloud.google.com/firestore/docs/quickstart-servers#create_a_in_native_mode_database) part of the quickstart guide.   As long as you will deploy Tomolink to the same GCP project as the Firestore instance, Tomolink will automatically be able to access Firestore to save data without additional authentication.  The [default configuration](../cmd/tomolink_defaults.yaml) uses a special `*detect-project-id*` value that instructs the Firestore client library to try to auto-detect your GCP project ID, and should work in most cases.
2. If you haven't already, enable the GCP [Cloud Run API](http://console.cloud.google.com/apis/library/run.googleapis.com?_ga=2.23508026.1852199750.1580108081-50451190.1521501879)
3. Start the build by running the following command from the repository root. It both builds the container and deploys it to Google Cloud Run. It will take a moment to complete. If you haven't yet enabled the Cloud Build API, the command will prompt you to do so, and provide you with a link.
```bash
gcloud builds submit --config cloudbuild.yaml .
```
4. You can view your deployment by going to [Cloud Run in the Google Cloud Console](https://console.cloud.google.com/run?enableapi=true&_ga=2.224795418.1852199750.1580108081-50451190.1521501879).
5. Finally, if you want to allow access to Tomolink from outside your GCP project, you'll need to [allow unauthenticated access](https://cloud.google.com/run/docs/authenticating/public).  **This is not recommended for production.**

## Setting up Continuous Deployment of Tomolink
You can use Cloud Build to automate builds and deployments to Cloud Run. See the guide on [Continuous Deployment from git using Cloud Build](https://cloud.google.com/run/docs/continuous-deployment-with-cloud-build).

## Building without deploying

To build the Tomolink Docker container and push it to Google Container Registry but NOT deploy it to Cloud Run using Cloud Build, modify the included [cloudbuild.yaml](../cloudbuild.yaml) file to remove the last step before submitting the `gcloud build submit --config cloudbuild.yaml .` command.
