
# Pulse
---
1. Created the Job model and helper functions
2. Create migrations for the Job
3. Enque( using repo interfaces, insert jobs in the DB)
4. Schedule, scheduler decides which tasks gets executed based on priority or scheduling.
5. Worker dequeues job from DB and checks the job
6. Based on the job type, the worker redirects the job to the appropriate handler and updates its state.