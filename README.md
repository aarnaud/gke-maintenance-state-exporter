# GKE/GCP maintenance state exporter

Prometheus exporter for GCP instance maintenance-event

base on https://cloud.google.com/compute/docs/storing-retrieving-metadata#query-events-sample-script

Possible value base on maintenance-event:

- NONE => 0
- MIGRATE_ON_HOST_MAINTENANCE => 1
- TERMINATE_ON_HOST_MAINTENANCE => 2
- any errors => -1

This exporter listen on port `9723` and the endpoint is `/metrics`

```
# HELP gcp_maintenance_state Report if a maintenance is planned on a GCE instance.
# TYPE gcp_maintenance_state gauge
gcp_maintenance_state{host="YOUR_INSTANCE_NAME"} 0
```

## TODO:

- [ ] Add maintenance metric on pod present on the node