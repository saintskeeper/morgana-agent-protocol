# Morgana Protocol Monitoring

This monitoring stack provides comprehensive observability for the Morgana
Protocol agent orchestration system using OpenTelemetry.

## Components

- **OpenTelemetry Collector**: Receives, processes, and exports telemetry data
- **Jaeger**: Distributed tracing for visualizing agent execution flows
- **Prometheus**: Metrics storage and querying
- **Grafana**: Dashboards for visualization

## Quick Start

1. Start the monitoring stack:

   ```bash
   docker-compose up -d
   ```

2. Access the UIs:

   - Grafana: http://localhost:3000 (admin/admin)
   - Jaeger: http://localhost:16686
   - Prometheus: http://localhost:9090

3. Run Morgana with OTLP exporter:
   ```bash
   ./dist/morgana --otel-exporter otlp -- --agent code-implementer --prompt "Write code"
   ```

## Telemetry Data

### Traces

The system creates detailed traces for:

- Overall execution (`morgana.execute`)
- Orchestrator operations (`orchestrator.sequential`, `orchestrator.parallel`)
- Individual agent executions (`agent.execute`)
- Sub-operations:
  - Agent validation
  - Prompt loading
  - Task execution
  - Semaphore waiting (for parallel execution)

### Attributes

Each span includes relevant attributes:

- `agent.type`: Type of agent (code-implementer, test-specialist, etc.)
- `agent.task_id`: Unique task identifier
- `agent.framework`: Always "morgana-protocol"
- `prompt.length`: Length of the combined prompt
- `result.success`: Whether execution succeeded
- `result.output_length`: Length of agent output
- `result.execution_time_ms`: Total execution time
- `task.count`: Number of tasks being orchestrated
- `task.index`: Index of current task

### Grafana Dashboard

The included dashboard shows:

1. **Agent Execution Latency**: P50 and P95 latencies by agent type
2. **Agent Execution Rate**: Current rate of agent executions
3. **Agent Type Distribution**: Pie chart of agent usage
4. **Success vs Failure Rate**: Time series of execution results

## Configuration

### Morgana CLI Options

- `--otel-exporter`: Choose exporter type
  - `stdout`: Print traces to console (default)
  - `otlp`: Send to OpenTelemetry Collector
  - `none`: Disable telemetry
- `--otel-endpoint`: OTLP endpoint (default: localhost:4317)

### Environment Variables

- `MORGANA_DEBUG=true`: Enable debug mode with verbose telemetry

## Example Usage

### Single Agent with Telemetry

```bash
./dist/morgana --otel-exporter otlp -- --agent code-implementer --prompt "Write a function"
```

### Parallel Agents with Telemetry

```bash
./dist/morgana --parallel --otel-exporter otlp -- \
  --agent code-implementer --prompt "Write code" \
  --agent test-specialist --prompt "Write tests"
```

### Debug Mode

```bash
MORGANA_DEBUG=true ./dist/morgana --otel-exporter stdout -- \
  --agent sprint-planner --prompt "Plan sprint"
```

## Troubleshooting

1. **No data in Grafana**:

   - Check that Prometheus datasource is configured
   - Verify OTLP collector is receiving data (check logs)
   - Ensure Morgana is using `--otel-exporter otlp`

2. **Traces not appearing in Jaeger**:

   - Check Jaeger UI at http://localhost:16686
   - Verify service name "morgana-protocol" is selected
   - Check collector logs for errors

3. **Connection refused**:
   - Ensure Docker containers are running: `docker-compose ps`
   - Check if ports are already in use
   - Verify firewall settings

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Morgana   │────▶│ OTEL         │────▶│   Jaeger    │
│   Protocol  │     │ Collector    │     │  (Traces)   │
└─────────────┘     └──────────────┘     └─────────────┘
                             │
                             ▼
                    ┌──────────────┐     ┌─────────────┐
                    │ Prometheus   │◀────│   Grafana   │
                    │  (Metrics)   │     │(Dashboards) │
                    └──────────────┘     └─────────────┘
```

## Future Enhancements

- Custom metrics for:
  - Prompt token counts
  - Agent-specific business metrics
  - Queue depths for parallel execution
- Alerting rules for:
  - High failure rates
  - Slow agent executions
  - Resource exhaustion
- Log aggregation with Loki
- Distributed tracing across Python bridge
