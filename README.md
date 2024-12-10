### Task Currency Exchange API
Objective
Develop a REST API that receives a JSON request with a currency pair (e.g., { "currency-pair": "USD-EUR" }) and retrieves the exchange rate for this pair.

Workflow
Concurrently request exchange rates from two distinct external services (Service A, Service B). Your API would call both of these services internally in the workflow.
Return the first response and disregard the second.
Return responses in the format: { "currency-pair": rate }, e.g., { "USD-EUR": 0.92 }.

### API documentation
The API documentation is available [here](./redoc-static.html).

```aiignore
POST /v1/exchange-rate


payload:

{
  "data": {
    "currency_pair": "USD-EUR"
  }
}


response:

{
  "data": {
    "USD-EUR": 0.92
  }
}
```

### Setup and run the application
Run the command in your terminal
``make env`` 
this would create your .env file

Next:
Run
``make run``
to start the application