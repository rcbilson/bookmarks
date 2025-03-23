import { Tabs } from '@chakra-ui/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Provider } from "@/components/ui/provider"

import FavoritePage from "./FavoritePage"
import RecentPage from "./RecentPage"
import SearchPage from "./SearchPage"
import AddBookmarkPage from "./AddBookmarkPage"

const queryClient = new QueryClient()

export default function App() {
  return (
    <Provider>
      <QueryClientProvider client={queryClient}>
        <Tabs.Root defaultValue="favorites" variant="line">
          <Tabs.List>
            <Tabs.Trigger value="favorites">
              Favorites
            </Tabs.Trigger>
            <Tabs.Trigger value="recent">
              Recent
            </Tabs.Trigger>
            <Tabs.Trigger value="search">
              Search
            </Tabs.Trigger>
            <Tabs.Trigger value="add">
              Add
            </Tabs.Trigger>
          </Tabs.List>
          <Tabs.Content value="favorites">
            <FavoritePage />
          </Tabs.Content>
          <Tabs.Content value="recent">
            <RecentPage />
          </Tabs.Content>
          <Tabs.Content value="search">
            <SearchPage />
          </Tabs.Content>
          <Tabs.Content value="add">
            <AddBookmarkPage />
          </Tabs.Content>
        </Tabs.Root>
      </QueryClientProvider>
    </Provider>
  )
}
