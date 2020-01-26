# mysql-event-demo

JUST A POC to demo for the team. Don't use the code in production.

### Problem Statement

When generating event from an API when stuffs are written for change,
the event operation might fail because of throttling, network issues etc

### Possible Solutions

1) Implement a retry policy that runs things in background
  ```
  Can work, but that means, we cannot restart instances which are 
  doing retries.
  ```
2) Two phase commits
  ```
  More complex
  Bad User Experience?  
  ```
3) Have a concurrent worker listening to events happening in
one API instance.
  ```
  Same problem as 1
  ``` 
  
4) Have a separate service reading from mysql binlog
  ```
  Not sure how multi server setup is going to work
  More moving part, but can be solved if run as a go process
  ```

### Presented POC (My View)

Having solution 2 & 4 will be more available and consistent system. (Throwing in buzz-words??)
But 2 is complex in terms of implementation.