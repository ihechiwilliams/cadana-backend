## Additional Task: Data Manipulation and Interfaces in Go

### Objective
Write a Go program that demonstrates data manipulation and interface usage, focusing on object-oriented practices.

### Tasks
1. **Unmarshal and Object Creation:**
    - Unmarshal a JSON file with 10 records in the format: `{ "id": "X", "personName": "Cadanaut X", "salary": { "value": "10", "currency": "USD" } }` into a struct called `Persons` for the array.

   ``` 
      types Persons struct {
         Data []Person `json:"data"`
      }
   ```

    - Design a `Person` object with appropriate fields and methods to encapsulate each row.

2. **Methods for Data Operations:**
    - Attach methods to the `Persons` struct to perform the following operations:
        - Sort the data array of `Person` objects by salary in ascending and descending order.
        - Group `Person` objects by salary currency into hash maps.
        - Filter `Person` objects by salary criteria in USD. Inject your API logic above as a dependecy to obtain the exchange rates to be able to filter all the different currencies in USD.
