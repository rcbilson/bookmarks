import { useCallback, useState } from 'react'
import { Tabs } from '@chakra-ui/react'
import { LuBookmarkPlus, LuStar, LuClock, LuSearch } from "react-icons/lu"
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Provider } from "@/components/ui/provider"

import FavoritePage from "./FavoritePage"
import RecentPage from "./RecentPage"
import SearchPage from "./SearchPage"
import AddBookmarkPage from "./AddBookmarkPage"
import ShareTarget from "./ShareTarget"

const queryClient = new QueryClient()

const isShareTarget = () => window.location.pathname.startsWith("/share-target")

export default function App() {
  const [sharing, setSharing] = useState(isShareTarget)
  const [tab, setTab] = useState("favorites")

  // Once a shared link has been added, drop the share-target URL and show
  // the recents list, where the new bookmark will appear.
  const onShareDone = useCallback(() => {
    window.history.replaceState(null, "", "/")
    queryClient.invalidateQueries({ queryKey: ['bookmarkList'] })
    setTab("recent")
    setSharing(false)
  }, [])

  return (
    <Provider>
      <QueryClientProvider client={queryClient}>
        {sharing ? <ShareTarget onDone={onShareDone} /> : (
        <Tabs.Root value={tab} onValueChange={(e) => setTab(e.value)} variant="line">
          <Tabs.List>
            <Tabs.Trigger value="favorites">
              <LuStar />
              Favorites
            </Tabs.Trigger>
            <Tabs.Trigger value="recent">
              <LuClock />
              Recent
            </Tabs.Trigger>
            <Tabs.Trigger value="search">
              <LuSearch />
              Search
            </Tabs.Trigger>
            <Tabs.Trigger value="add">
              <LuBookmarkPlus />
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
        )}
      </QueryClientProvider>
    </Provider>
  )
}
