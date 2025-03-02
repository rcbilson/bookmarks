// Supports weights 100-900
//import '@fontsource-variable/outfit';

import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'

import ErrorPage from "./ErrorPage.jsx";
import RecentPage from "./RecentPage.tsx";
import FavoritePage from "./FavoritePage.tsx";
import SearchPage from "./SearchPage.tsx";

const router = createBrowserRouter([
  {
    path: "/",
    element: <FavoritePage />,
    errorElement: <ErrorPage />,
  },
  {
    path: "/recent",
    element: <RecentPage />,
  },
  {
    path: "/favorite",
    element: <FavoritePage />,
  },
  {
    path: "/search",
    element: <SearchPage />
  }
]);

const queryClient = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  )
}
