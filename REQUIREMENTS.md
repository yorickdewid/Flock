## Application requirements

- RFQ1: Create, remove, list queues
- RFQ2: Submit new jobs per queue and return job UUID
- RFQ3: Fetch next job in queue by UUID
- RFQ4: Purge all jobs in queue
- RFJ5: Request job status by UUID
- RFJ6: Cancel job by UUID
- RFT7: Request task status by UUID
- RFT8: Cancel task by UUID

## Quality requirements

- RN1: Fast API response
- RN2: Backwards compatible for at least 2 API versions
- RN3: Fatal errors MUST return response
