# ns

The ns project is a simple task processing mock example. 

M producers create tasks and store it into channel, N consumers handle these tasks and consume it.

To stop all workers we should have another worker manager which handle a stop signal and stop all workers and hence we stick to the rule of the channel closing principle.
