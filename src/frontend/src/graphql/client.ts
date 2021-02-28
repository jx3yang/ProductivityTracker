import { ApolloClient, InMemoryCache } from "@apollo/client"

const host = "localhost"
const port = "8080"
const graphqlPath = "query"

// https://www.apollographql.com/docs/react/caching/cache-field-behavior/#merging-non-normalized-objects
const cache = new InMemoryCache({
  typePolicies: {
    List: {
      fields: {
        cards: {
          merge(_, incoming) {
            return incoming;
          }
        }
      }
    }
  }
})

export const client = new ApolloClient({
  uri: `http://${host}:${port}/${graphqlPath}`,
  cache: cache,
});
