# orchestrator

#### License

`orchestrator` is free and open sourced under the [Apache 2.0 license](LICENSE).

#### SDK

Add Cascade or Double Master switchover

useage:

`orchestrator-client -c cascade-master-takeover-auto -i slave_inst:port -d dest_inst:port`

Add precheck tools, check list:
- cascade or double master replication topology
- orchestrator nodes ssh to topology check
- meta.cluster table check
- orchestrator dect mysql account connection check
- dect mysql account privileges check
- hostname dns check
- metadata check
