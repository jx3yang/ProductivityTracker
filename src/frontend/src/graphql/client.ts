import { ApolloClient, InMemoryCache } from "@apollo/client"

const host = "localhost"
const port = "8080"
const graphqlPath = "query"

export const client = new ApolloClient({
  uri: `http://${host}:${port}/${graphqlPath}`,
  cache: new InMemoryCache()
});
