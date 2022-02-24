# Colosseum Near-Real-Time RIC

This is a part of the [OpenRAN Gym](https://openrangym.com) project. It is minimal version of the O-RAN Software Community near-real-time RIC (Bronze release) adapted and extended to work on the [Colosseum](https://www.colosseum.net/) wireless network emulator.
The scripts in this repository will start a minimal near-real-time RIC in the form of Docker containers (namely, `db`, `e2mgr`, `e2rtmansim`, `e2term`).
The repository also features a sample xApp, which connects to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2).

If you use this software, please reference the following paper: 

> M. Polese, L. Bonati, S. D'Oro, S. Basagni, T. Melodia, "ColoRAN: Design and Testing of Open RAN Intelligence on Large-scale Experimental Platforms," arXiv 2112.09559 [cs.NI], December 2021. [bibtex](https://ece.northeastern.edu/wineslab/wines_bibtex/polese2021coloran.txt) [pdf](https://arxiv.org/pdf/2112.09559.pdf)

This work was partially supported by the U.S. National Science Foundation under Grants CNS-1923789 and NSF CNS-1925601, and the U.S. Office of Naval Research under Grant N00014-20-1-2132.

## Structure

This repository is organized as follows

```
root 
|
└──setup-scripts
|   |
|   └──import-base-images.sh
|   |
|   └──setup-lib.sh
|   |
|   └──setup-ric.sh
|   |
|   └──setup-sample-xapp.sh
|   |
|   └──start-ric-arena.sh
|   |
|   └──start-xapp.sh
|   
└──setup
|  |
|  └──dbaas
|  |
|  └──e2
|  |
|  └──sample-xapp
|  |
|  └──xapp-sm-connector
```

### Quick start

We provide a Colosseum LXC container that contains this repository, its prerequisites, and base docker images. The container `coloran-near-rt-ric` can be found among the images available for Colosseum users. The default username and password are `root` and `ChangeMe`.

From the `setup-scripts` directory:
- Build, configure, and start the near-real-time RIC Docker containers: `./setup-ric.sh ric-network-interface`
- Connect the RAN node through the E2 termination as explained [here](https://github.com/wineslab/colosseum-scope-e2)
- Get the gNB ID (see section below) and replace it in the `start-xapp.sh` script
- Configure and start the xApp: `./start-xapp.sh`

### setup-scripts directory

The `setup-scripts` directory contains scripts to initialize the near-real-time RIC on Colosseum.
- `import-base-images.sh`: script to import the base Docker images needed to build the RIC Docker containers. These images are provided as part of the `coloran-near-rt-ric` Colosseum LXC container.
- `setup-lib.sh`: contains the IP addresses and ports used by the Docker containers of this repository. This script has been adapted from [here](https://gitlab.flux.utah.edu/johnsond/ric-profile/-/blob/master/setup-lib.sh)
- `setup-ric.sh`: script to build, configure, and start the near-real-time RIC containers of this repository (namely, `db`, `e2mgr`, `e2rtmansim`, `e2term`). The network interface the RIC listens to for connections (e.g., the `col0` interface in Colosseum) is passed as argument. This script has been adapted from [here](https://gitlab.flux.utah.edu/johnsond/ric-profile/-/blob/master/setup-ric.sh)
- `start-ric-arena.sh`: script to start the near-real-time RIC on external testbeds, e.g., on the [Arena platform](https://ece.northeastern.edu/wineslab/arena.php)
- `setup-sample-xapp.sh`: script to setup a sample xApp Docker container. This xApp is capable of connecting to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2). Custom or standard-compliant service models can be implemented on top of the RAN E2 termination and the sample xApp, as done for example [in these]() [papers](https://ece.northeastern.edu/wineslab/papers/bonati2021intelligence.pdf)
- `start-xapp.sh`: script to configure and start the sample xApp. The ID of the gNB targeted by the xApp needs to be provided in the script, as discussed below

### setup directory

This directory contains the implementations of the near-real-time RIC Docker container initialized through the scripts in the [`setup-scripts`](setup-scripts) directory.
- `dbaas`: implementation of a Redis database (`db`) container
- `e2`: implementation of the E2 termination (`e2term`) container
- `e2mgr`: implementation of the E2 manager (`e2mgr`) and of the routing manager simulator (`e2rtmansim`) container
- `sample-xapp`, `xapp-sm-connector`: implementation of the sample xApp provided in this repository and components to connect to the near-real-time RIC and [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment

These components are adapted from the [O-RAN Software Community RIC platform (Bronze)](https://github.com/o-ran-sc), which we extended to support the Colosseum environment, concurrent connections from multiple base stations and xApps, and to provide improved support for encoding, decoding and routing of control messages.

### Getting the gNB ID

An easy way to derive the gNB ID is the following. From the `setup-scripts` directory:
- Start the near-real-time RIC Docker containers: `./setup-ric.sh`
- Read the logs of the `e2term` container: `docker logs e2term -f`
- Connect the RAN node through the E2 termination as explained [here](https://github.com/wineslab/colosseum-scope-e2)
- The RAN node should connect to the near-real-time RIC and the gNB ID should appear in the `e2term` logs. In the example below, the gNB ID is `gnb:311-048-01000501`

  ```
  {"ts":1639008174427,"crit":"DEBUG","id":"E2Terminator","mdc":{"thread id":"139898725332736"},"msg":"After processing message and sent to rmr for : gnb:311-048-01000501, Read time is : 0 seconds, 1044889 nanoseconds"}
  ```

### Using the provided sample xApp

The sample xApp provided in this repository connects to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2).
After the near-real-time RIC has successfully started, the DU connected to it, and the xApp has been properly configured and started (see "Quick start" section above):
- Enter the xApp docker container (named `sample-xapp-24` by default): `docker exec -it sample-xapp-24`
- Move to the `/home/sample-xapp` directory inside the Docker container: `cd /home/sample-xapp`
- Run the xApp logic: `./run_xapp.sh`. This script will open a socket between the sample Python script in the `sample-xapp` directory (which by defaults prints the data received from the RAN node) and the service model connector of the `xapp-sm-connector` directory, which performs ASN.1 encoding and decoding of E2AP messages. Then, the xApp will subscribe to the RAN node specified at container startup time through the gNB ID, and receive a RIC Indication Message with a data report from the RAN node with the periodicity of 250 ms.
