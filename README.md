# oras-operator

Deploy an ORAS registry (cache for workflow or experiment artifacts) as a service.

![docs/assets/design/im-quite-the-catch.png](docs/assets/design/im-quite-the-catch.png)

- Read our 🌟️ [Documentation](https://converged-computing.github.io/oras-operator) 🌟️ to get started.

## TODO:

- test out setup with scripts, merge when basics are working
- create docs and automated builds for containers
- watch job (and then edit template) insteaed of pod (or make a variable, better?)
- test with a simple dag (maybe snakemake kueue executor)
- multiple pods for registry and using secret / shared storage use case


## License

HPCIC DevTools is distributed under the terms of the MIT license.
All new contributions must be made under this license.

See [LICENSE](https://github.com/converged-computing/cloud-select/blob/main/LICENSE),
[COPYRIGHT](https://github.com/converged-computing/cloud-select/blob/main/COPYRIGHT), and
[NOTICE](https://github.com/converged-computing/cloud-select/blob/main/NOTICE) for details.

SPDX-License-Identifier: (MIT)

LLNL-CODE- 842614
