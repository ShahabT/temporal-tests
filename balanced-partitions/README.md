Balanced Partition Test
=======================
This test is a series of experiments aimed for verifying the TQ forwarding logic to make sure:
- Tasks and pollers are fairly distributed among the partitions
- Backlog is accumulated and processed equally in all partitions
- When there are 'too many' or 'just enough' workers, the tasks should primarily sync match
- Forwarding happens as needed to prevent pollers being stuck in a partition without task
- When a partition has no task to dispatch, it does not black hole for pollers (even if that is a root partition)


Instructions
------------

1. Run the server
2. Run the worker like bellow. Replace `<exp_name>` with an experiment name from `balanced-partitions/experiments.go`.
    ```
    $ go run start/main.go <exp_name>
    ```
3. Use the dashboard at `graphana-dashboard.json` to monitor the metrics and verify behavior
