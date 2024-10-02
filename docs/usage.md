# Usage Guide for ICP-EVM Proxy

This guide provides instructions on how to interact with the ICP-EVM Proxy system after it has been set up and deployed.

## Interacting with the Counter Canister

1. **Increment the Counter**

   Use the following command to increment the counter:

   yyyshell
   dfx canister call counter_canister increment
   yyy

2. **Get the Current Counter Value**

   To retrieve the current value of the counter:

   yyyshell
   dfx canister call counter_canister get
   yyy

   Both of these operations will generate events that are logged by the Logger Canister.

## Querying Logs via EVM Adapter Proxy

You can interact with the EVM Adapter Proxy using standard Ethereum JSON-RPC calls. Here are some example queries:

1. **Get the Latest Block Number**

   yyyshell
   curl -X POST -H "Content-Type: application/json" \
   --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
   <http://localhost:3030/rpc/v1>
   yyy

2. **Get Logs**

   To retrieve logs for a specific block range:

   yyyshell
   curl -X POST -H "Content-Type: application/json" \
   --data '{
     "jsonrpc":"2.0",
     "method":"eth_getLogs",
     "params":[{"fromBlock":"0x0","toBlock":"latest"}],
     "id":1
   }' \
   <http://localhost:3030/rpc/v1>
   yyy

## Querying Data via SubQuery Indexer

The SubQuery Indexer provides a GraphQL endpoint for querying indexed data. You can access the GraphQL playground at `http://localhost:3000`.

Here's an example query to get the most recent logs:

```graphql
query {
  logs(
    first: 5,
    orderBy: [BLOCK_NUMBER_DESC]
  ) {
    nodes {
      id
      address
      topics
      data
      blockNumber
      transactionHash
      transactionIndex
      blockHash
      logIndex
      removed
    }
    totalCount
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
  }
}
```

You can modify this query to filter logs based on specific criteria or to retrieve different fields.

### Accessing the Database Directly

You also have the option to query the database directly. The SubQuery Indexer typically uses PostgreSQL as its database. To access it:

1. Connect to the database using your preferred PostgreSQL client.
2. The default connection details are usually:
   - Host: localhost
   - Port: 5432
   - Database: postgres
   - User: postgres
   - Password: postgres (unless changed during setup)

3. Once connected, you can query the `logs` table directly. For example:

```sql
SELECT * FROM <schema_name>.logs;
```

This SQL query is equivalent to the GraphQL query shown above. You can modify it to filter or retrieve data as needed.

## Troubleshooting

### GraphQL Error with public.subqueries Table

If you encounter an error related to the `public.subqueries` table when trying to use the GraphQL endpoint, you may need to create this table manually and insert a record. Here's how to do it:

1. Connect to your PostgreSQL database using your preferred client.

2. Create the `public.subqueries` table if it doesn't exist:

   ```sql
   CREATE TABLE IF NOT EXISTS public.subqueries (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       version VARCHAR(255) NOT NULL
   );
   ```

3. Insert a record for your application:

   ```sql
   INSERT INTO public.subqueries (name, version)
   VALUES ('poc-subquery-icp-icrc3', '0.0.1');
   ```

   Note: Replace 'poc-subquery-icp-icrc3' with the actual name of your application if different.

4. After executing these SQL commands, try accessing the GraphQL endpoint again. The error should be resolved.
