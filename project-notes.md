# Monitoring: Prometheus and Grafana
Added Promtheus and Grafana. Had to mess around with the configurations and learn how to point a Docker container to a yaml file to properly configure Prometheus to scrape the correct endpoint. 

Additionally learned about Docker container ports and how to configure containers to point outside of them to host ports. Will likely need to reconfigure this in the future.

# ToDo
- status updates for batch jobs
- handle failed jobs
- set up error handling for the workers
- add new monitoring metrics for the worker service
    - USE method