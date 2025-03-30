// A react component that has an editable text area for a bookmark url
// next to a button with a refresh icon. When the button is clicked,
// the bookmark url is fetched and the text area below the url is updated
// with the bookmark contents.
import React from "react";
import axios from "axios";
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { HStack, VStack, Box } from "@chakra-ui/react"
import { LuStar } from "react-icons/lu";

type BookmarkEntry = {
  title: string;
  url: string;
  isFavorite: boolean;
}

interface Props {
  queryPath: string;
}

const BookmarkQuery: React.FC<Props> = ({ queryPath }: Props) => {
  const queryClient = useQueryClient();

  const fetchQuery = (queryPath: string) => {
    return async () => {
      console.log("fetching " + queryPath);
      const response = await axios.get<Array<BookmarkEntry>>(queryPath);
      return response.data;
    };
  };

  const { isError, data, error } = useQuery({
    queryKey: ['bookmarkList', queryPath],
    queryFn: fetchQuery(queryPath),
  });
  const recents = data;

  const handleBookmarkClick = (url: string) => {
    return () => {
      const encodedUrl = encodeURIComponent(url);
      axios.post("/api/hit?url=" + encodedUrl);
      window.open(url, "_blank");
    }
  }

  const handleStarClick = (url: string, isFavorite: boolean) => {
    return () => {
      const encodedUrl = encodeURIComponent(url);
      axios.post(`/api/setFavorite?url=${encodedUrl}&isFavorite=${isFavorite}`).then(
        (_) => {
          queryClient.invalidateQueries({ queryKey: ['bookmarkList'] });
        }
      );
    }
  }

  if (isError) {
    return <div>An error occurred: {error.message}</div>
  }

  return (
    <div id="bookmarkList">
      {recents && recents.map((recent) =>
        <HStack>
          <Box w="20px">
            <LuStar onClick={handleStarClick(recent.url, !recent.isFavorite)} color={recent.isFavorite ? "gold" : "gray"} size={20} />
          </Box>
          <VStack align="left" spaceY={0} >
            <div className="bookmarkEntry" key={recent.url} onClick={handleBookmarkClick(recent.url)}>
              <div className="title">{recent.title}</div>
              <div className="url">{new URL(recent.url).hostname}</div>
            </div>
          </VStack>
        </HStack>
      )}
    </div>
  );
};

export default BookmarkQuery;
