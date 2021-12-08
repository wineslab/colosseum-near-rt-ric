# Colosseum Near-Real-Time RIC

This is a minimal implementation of the O-RAN Software Community near-real-time RIC (Bronze release) adapted and extended to work on the [Colosseum](https://www.colosseum.net/) wireless network emulator.
The scripts in this repository will start a minimal near-real-time RIC in the form of Docker containers (namely, `db`, `e2mgr`, `e2rtmansim`, `e2term`).
A sample xApp, which connects to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2) is also provided.

If you use this software, please reference the following paper:

XXX

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
|  └──xapp-bs-connector
```

### Quick start

From the `setup-scripts` directory:
- Build, configure, and start the near-real-time RIC Docker containers: `./setup-ric.sh`
- Connect the RAN node through the E2 termination as explained [here](https://github.com/wineslab/colosseum-scope-e2)
- Get the gNB ID (see section below) and replace it in the `start-xapp.sh` script
- Configure and start the xApp: `./start-xapp.sh`

### setup-scripts directory

The `setup-scripts` directory contains scripts to initialize the near-real-time RIC on Colosseum.
- `import-base-images.sh`: script to import the base Docker images needed to build the RIC Docker containers. These images are provided as part of the XXX Colosseum LXC container.
- `setup-lib.sh`: contains the IP addresses and ports used by the Docker containers of this repository. This script has been adapted from [here](https://gitlab.flux.utah.edu/johnsond/ric-profile/-/blob/master/setup-lib.sh)
- `setup-ric.sh`: script to build, configure, and start the near-real-time RIC containers of this repository (namely, `db`, `e2mgr`, `e2rtmansim`, `e2term`). This script has been adapted from [here](https://gitlab.flux.utah.edu/johnsond/ric-profile/-/blob/master/setup-ric.sh)
- `start-ric-arena.sh`: script to start the near-real-time RIC on external testbeds, e.g., on the [Arena platform](https://ece.northeastern.edu/wineslab/arena.php)
- `setup-sample-xapp.sh`: script to setup a sample xApp Docker container. This xApp is capable of connecting to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2)
- `start-xapp.sh`: script to configure and start the sample xApp. The ID of the gNB targeted by the xApp needs to be provided in the script

### setup directory

This directory contains the implementations of the near-real-time RIC Docker container initialized through the scripts in the [`setup-scripts`](setup-scripts) directory.
- `dbaas`: implementation of a Redis database (`db`) container
- `e2`: implementation of the E2 termination (`e2term`) container
- `e2mgr`: implementation of the E2 manager (`e2mgr`) and of the routing manager simulator (`e2rtmansim`) container
- `sample-xapp`, `xapp-bs-connector`: implementation of the sample xApp provided in this repository and components to connect to the near-real-time RIC and [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment

### Getting the gNB ID

An easy way to derive the gNB ID is the following. From the `setup-scripts` directory:
- Start the near-real-time RIC Docker containers: `./setup-ric.sh`
- Read the logs of the `e2term` container: `docker logs e2term -f`
- Connect the RAN node through the E2 termination as explained [here](https://github.com/wineslab/colosseum-scope-e2)
- The RAN node should connect to the near-real-time RIC and the gNB ID should appear in the `e2term` logs (see example below)

```
TODO insert sample gNB ID from e2term logs
```

### Using the provided sample xApp

The sample xApp provided in this repository connects to the [SCOPE](https://github.com/wineslab/colosseum-scope) RAN environment through the following [E2 termination](https://github.com/wineslab/colosseum-scope-e2).
After the near-real-time RIC has successfully started, the DU connected to it, and the xApp has been properly configured and started (see "Quick start" section above):
- Enter the xApp docker container (named `sample-xapp-24` by default): `docker exec -it sample-xapp-24`
- Move to the right directory inside the Docker container: `cd /home/sample-xapp`
- Start the xApp: `./run_xapp.sh`. This script will open a socket between the sample script in the `sample-xapp` directory (which by defaults prints the data received from the RAN node) and the base station connector of the `xapp-bs-connector` directory. Then, the xApp will subscribe to the RAN node specified at contianer startup time through the gNB ID, and send a RIC Indication Message that triggers a data report from the RAN node with the periodicity of 1 second.
