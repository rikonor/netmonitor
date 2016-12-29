netmonitor
---

docker run --rm -it -p 7474:7474 -p 7687:7687 neo4j

Chrome extension
  - Intercepts outgoing requests and forwards their details to a collector
    - Before/After they come back? Do we have status, etc?
  - Can we put this patched chrome in a docker container?
  - Can we provide the extension with the collector address via env var?
  - Can we provide the collector address when creating the docker image?
  - What about websockets?

Worker
  - Contains chrome with the extension
  - Do we need a worker? Maybe just the collector starting a chrome docker container with a suitable cmd?
    That said, we then have much less control
    Also, the worker can act as a local collector and forward whatever the local chrome sends it
    then its trivial to pass along configuration details

Collector
  - Collects the request details from the Chrome extension
  - Each request detail should contain also which website made the request

Questions
  - Do different devices invoke different requests? m.wikipedia.com, etc
  - What software and what software version do they sites load? (jQuery v0.8, etc), old software could be vulnerable.
  - Related vendors (site calls 3rd party site, etc)
  - Can we graph this data? Sites using certain tech, sites connected to each other (probably uni-directional for the most part)
  - What about logins to sites? That would require selenium. after login sites might load different assets
