# clients

This repository is for holding all client codes to micro services.

Please implement the client code of future services in this repo, \
and the process of moving existing client into this repo is being carried on. \


# code structure

Please follow the elaborated repository structure when designing new components

```
             ┌-----→  lib-common  ←--┐
             |             ↑         |
          services ( → ) clients     |
                           ↑         |
                      applications---┘
```

- lib-common
    - contains types that are shared across multiple components
    - contains general functions used in multiple components
- services
    - depends on lib-common for shared types
    - depends on client repository if inter-service communication is required
- clients
    - depends on lib-common for shared types
    - should not depend on service repos
- applications
    - depends on clients for accessing services
    - depends on lib-common for shared types

# summary

services | client code status | sharing custom types
---------|--------------------|------------
cms      | implemented  |  no
aas      | implemented  |  yes