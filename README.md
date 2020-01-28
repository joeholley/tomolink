# Tomolink
Tomolink is an open source friend/block/relationship service for games. It allows you to store, retrieve, and delete arbitrary directed relationships at scale. 

Tomolink can track a user's `friends`, the `influencers` they follow, and the other users they have chosen to `block` in its default configuration, but by arbitrary, we mean it: you can choose _nearly any string_ to represent a relationship you want to track.  Do you want to track that one user `supports` another, or has a certain level of `distrust`, or has a `friendRequestPending`?  Tomolink can help! Not only does it offer flexibility in the type of relationships you want to track, Tomolink stores all relationships with an associated integer _score_. This allows you to track the significance of the relationship in addition to it's existence, and enables many exciting possibilities!  For more details on how to use Tomolink effectively, see the [User Guide](docs/userguide.md).

Note that Tomolink is architected on top of Google Cloud Platform managed products, to minimize any administration and management overhead.  For most use cases, setting up Tomolink is as simple as signing up for GCP, creating a project, enabling the products used by Tomolink, and running a couple of commands to start the service. It is designed to be cost-efficient for nearly all pertinent commercial game use cases.

## Contributing to Tomolink

Tomolink is in active development and we would love your contribution after
we've stabilized and cut a release! Please read the [contributing
guide](CONTRIBUTING.md) for guidelines on contributing to
Tomolink.

The [Tomolink Development guide](docs/development.md) has detailed instructions
on getting the source code, making changes, testing and submitting a pull request
to Tomolink.

## Disclaimer

**This software is currently alpha, and subject to change.  It may contain bugs and is used at your own risk; the authors assume no liability!**

## Support

Please file a github issue if you think you've found a bug. 

## Code of Conduct

Participation in this project comes under the [Contributor Covenant Code of Conduct](code-of-conduct.md)

## License

Apache 2.0
