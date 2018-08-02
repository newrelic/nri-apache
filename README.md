# New Relic Infrastructure Integration for Apache

New Relic Infrastructure Integration for Apache captures critical performance
metrics and inventory reported by Apache web server.

Inventory data is obtained using `httpd` command in RedHat family distributions
and `apache2ctl` command in Debian family distributions.
Metrics data is obtained doing HTTP requests to `/server-status` endpoint, provided by
`mod_status` Apache module.

## Usage
This is the description about how to run the Apache Integration with New Relic
Infrastructure agent, so it is required to have the agent installed
(see
[agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

In order to use the Apache Integration it is required to configure
`apache-config.yml.sample` file. Firstly, rename the file to
`apache-config.yml`. Then, depending on your needs, specify all instances that
you want to monitor with correct arguments.

The integration assumes your `/server-status` is reachable though the next URL:
```
http://127.0.0.1/server-status?auto
```

You can change this URL by modifying the `status_url` instance parameter of your `apache-config.yml` file.

If the `/server-status` is reachable only though an HTTPS connection which uses a certificate that is signed by your own
Certificate Authority (e.g. for internal use or development purposes), you will need to make sure **one** of the next alternative
actions have been taken:

* *Option 1*: Install your Certificate Authority in your host. The Apache integration will look for it by default in the host's root
  bundle for Certificate Authorities.
* *Option 2*: Update the `apache-config.yml` to have at least one of the next instance parameters configured: `ca_bundle_file` or
  `ca_bundle_dir`, whose values must be, respectively, the absolute paths to your alternative Certificate Authorities' bundle
  file (or the directory where they are).
* *Option 3*: Add to the [Infrastructure Agent configuration file](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/configuration/configure-infrastructure-agent))
  the `ca_bundle_file` or `ca_bundle_dir` property (which are analogue to the properties explained in the previous
  bullet), and then passthrough them to the integrations by adding: `passthrough_environment: CA_BUNDLE_FILE` or
  `passthrough_environment: CA_BUNDLE_DIR` to the `infra-agent.yml` configuration file.

## Integration development usage

Assuming that you have the source code and Go tool installed you can build and run the Apache Integration locally.
* After cloning this repository, go to the directory of the Apache Integration and build it
```bash
$ make
```
* The command above will execute the tests for the Apache Integration and build an executable file called `nr-apache` under `bin` directory. Run `nr-apache`:
```bash
$ ./bin/nr-apache
```
* If you want to know more about usage of `./bin/nr-apache` check
```bash
$ ./bin/nr-apache -help
```

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.

## Contributing Code

We welcome code contributions (in the form of pull requests) from our user
community. Before submitting a pull request please review [these guidelines](https://github.com/newrelic/nri-apache/blob/master/CONTRIBUTING.md).

Following these helps us efficiently review and incorporate your contribution
and avoid breaking your code with future changes to the agent.

## Custom Integrations

To extend your monitoring solution with custom metrics, we offer the Integrations
Golang SDK which can be found on [github](https://github.com/newrelic/infra-integrations-sdk).

Refer to [our docs site](https://docs.newrelic.com/docs/infrastructure/integrations-sdk/get-started/intro-infrastructure-integrations-sdk)
to get help on how to build your custom integrations.

## Support

You can find more detailed documentation [on our website](http://newrelic.com/docs),
and specifically in the [Infrastructure category](https://docs.newrelic.com/docs/infrastructure).

If you can't find what you're looking for there, reach out to us on our [support
site](http://support.newrelic.com/) or our [community forum](http://forum.newrelic.com)
and we'll be happy to help you.

Find a bug? Contact us via [support.newrelic.com](http://support.newrelic.com/),
or email support@newrelic.com.

New Relic, Inc.




