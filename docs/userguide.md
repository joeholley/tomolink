# Tomolink User Guide

At its most basic level, Tomolink is an HTTP API in front of a NoSQL document database. All the heavy lifting of storing user relationships is done by the database. The API focuses on performing basic validation of input parameters.  

The API provides 6 endpoints:

1) `http://<your_domain>/createRelationship` with parameters in the request JSON to create a relationship
1) `http://<your_domain>/updateRelationship` with parameters in the request JSON to update a relationship
1) `http://<your_domain>/deleteRelationship` with parameters in the request JSON to delete a relationship
1) `/users/<uuidsource>` to retrieve all relationships for the provided user ID.  
1) `/users/<uuidsource>/<relationship>` to retrieve all relationships of the given type for the provided user ID. 
1) `/users/<uuidsource>/<relationship>/<uuidtarget>` to retrieve the value of one relationship from the provided user ID to the target user ID. 

## Updating Configuration

Tomolink accepts a YAML config file called [tomolink_defaults.yaml](../cmd/tomolink_defaults.yaml). All of the values in the config file can also be overridden by environment variable. To do so, set an environment variable with the same name as the _dot notation of the YAML config parameter_, with all upper-case letters, and underscores in place of periods. For example, the config parameters for setting up the HTTP API in the YAML file look like this:
```yaml
http:
    port: 8080         # Port to serve on
    gracefulwait: 15   # Seconds to wait for requests to finish if graceful shutdown of server is requested
    request:
        readLimit: 500 # Limit the size of incoming requests to something sensible, abuse prevention measure
```

They could be overwritten by setting the following environment variables:
```bash
HTTP_PORT=8000
HTTP_GRACEFULWAIT=20
HTTP_REQUEST_READLIMIT=1200
```

### Setting a default config

If you want to change the default configuration for your Tomolink deployments (for example, to [specify relationship types](#choosing-relationships)), you can make changes to [tomolink_defaults.yaml](../cmd/tomolink_defaults.yaml) and rebuild the Tomolink docker container](development.md#building-tomolink).  

### Config on Cloud Run

If you want to specify some configuration overrides that differ from the default config file on a per-deployment basis, Cloud Run makes this quite easy.  Just [set the desired environment variables](https://cloud.google.com/run/docs/configuring/environment-variables) in the Cloud Run console.

### Confirming config settings

If you want to verify that your overridden config settings are being used by Tomolink, look in the logs.  Tomolink outputs a line for every config parameter override it processes on startup.

## Choosing the Relationships

Tomolink can track a user's `friends`, the `influencers` they follow, and the other users they have chosen to `block` in its default configuration, but by arbitrary, we mean it: you can choose _nearly any string_ to represent a relationship you want to track.  Do you want to track that one user `supports` another, or has a certain level of `distrust`, or has a `friendRequestPending`?  Tomolink can help! Not only does it offer flexibility in the type of relationships you want to track, Tomolink stores all relationships with an associated integer _score_. This allows you to track the significance of the relationship in addition to it's existence, and enables many exciting possibilities!  

### Strict vs non-strict
Tomolink has the option to turn on or off **"strict"** relationships.  

With strict relationships **disabled**, any create or update API call will try to save the relationship specifed in the provided [request JSON](#sending-input-parameters-in-the-json-body). 


When strict relationships are **enabled**, Tomolink will only accept API calls that specify one of the (up to) ten relationships defined in the [configuration](#updating-configuration).  All configured relationships have a [`name`](#relationship-names) 

### Relationship Names

Most strings are fine, but avoid using periods: refer to the [Field Names](https://cloud.google.com/firestore/docs/best-practices#field_names) section of the Firestore Best Practices documentation if you want to learn more about the limitations of what strings you can use for relationship types.

## Using Scores


## Relationship parameters
When sending an update to the Tomolink service, you must always specify the source and target user IDs, the relationship, the relationship _direction_, and the change to the [score](#using-scores) (the _delta_). Some illustrative examples:

User `d7e86e48-f8b5-48de-ad22-13c944b1d437` (who we'll call 'Dee' for short) wants to add user `f170dba6-c825-4fef-92f8-324351cd4908` ('Eff') to their friends list, with a score of '100' and vice-versa. In this case:

* the **source user ID** is Dee's: `d7e86e48-f8b5-48de-ad22-13c944b1d437`  
* the **target user ID** is Eff's: `f170dba6-c825-4fef-92f8-324351cd4908`
* the **relationship** is `friends`
* Dee is adding Eff, and vice-versa. The **direction** of this relationship is `mutual`.
* The score (called the **delta**) is `100` 

Now imagine Dee has a rough day Dee and Eff have a disagreement. Dee decides to block Eff, but Eff doesn't want to block Dee just yet.  This would only update the block relationship in one direction:

* the **source user ID** is still Dee's: `d7e86e48-f8b5-48de-ad22-13c944b1d437`  
* the **target user ID** is still Eff's: `f170dba6-c825-4fef-92f8-324351cd4908`
* the **relationship** is `blocks`
* Dee is adding Eff, but not vice-versa! The **direction** of this relationship is `single`.
* The score (called the **delta**) is `1`. 

### Sending input parameters in the JSON body

These three API calls expect you to specify the parameters of your request in the body using JSON:
* `/createRelationship`
* `/updateRelationship`
* `/deleteRelationship`

The request JSON body has the following format:
```json
{
    "uuidsource": "d7e86e48-f8b5-48de-ad22-13c944b1d437",
    "uuidtarget": "f170dba6-c825-4fef-92f8-324351cd4908",
    "relationship": "friends",
    "delta": 1000,
    "direction": "mutual",
}
```

The fields are as follows:
* `relationship` is the name of the relationship to act on
* `uuidsource` is the user ID whose relationship is being updated
* `uuidtarget` is the user ID 
* `direction` should be set to `single` to update the `uuidsource` relationship to `uuidtarget`. To also update `uuidtarget`'s relationship to `uuidsource`, set it to `mutual`. 
* `delta` is the change to the relationship score. In the case of `createRelationship`, the relationship score will be set to this value.  In `updateRelationship`, this value will be added to the current value (to subtract, provide a negative `delta` value). If this field is provided to `deleteRelationship`, it is ignored.

### Sending input parameters in the URI
These three API calls expect you to specify the parameters of your request in the request URI: 
* `/users/<uuidsource>`
* `/users/<uuidsource>/<relationship>`
* `/users/<uuidsource>/<relationship>/<uuidtarget>`
