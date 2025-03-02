// A react component that has an editable text area for a bookmark url
// next to a button with a refresh icon. When the button is clicked,
// the bookmark url is fetched and the text area below the url is updated
// with the bookmark contents.
import React from "react";
import { useNavigate } from 'react-router-dom';
import axios from "axios";
import { useQuery } from '@tanstack/react-query'

type BookmarkEntry = {
  title: string;
  url: string;
}

interface Props {
  queryPath: string;
}

const BookmarkQuery: React.FC<Props> = ({queryPath}: Props) => {
  const navigate = useNavigate();
 
  const fetchQuery = (queryPath: string) => {
    return async () => {
      console.log("fetching " + queryPath);
      const response = await axios.get<Array<BookmarkEntry>>(queryPath);
      return response.data;
    };
  };

  const {isError, data, error} = useQuery({
    queryKey: ['bookmarkList', queryPath],
    queryFn: fetchQuery(queryPath),
  });
  const recents = data;

  const handleBookmarkClick = (url: string) => {
    return () => {
      const encodedUrl = encodeURIComponent(url);
      axios.post("/api/hit?url=" + encodedUrl);
      navigate("/show/" + encodedUrl);
    }
  }
 
  if (isError) {
    return <div>An error occurred: {error.message}</div>
  }

  return (
    <div id="bookmarkList">
      {recents && recents.map((recent) =>
        <div className="bookmarkEntry" key={recent.url} onClick={handleBookmarkClick(recent.url)}>
          <div className="title">{recent.title}</div>
          <div className="url">{new URL(recent.url).hostname}</div>
        </div>
      )}
    </div>
  );
};

export default BookmarkQuery;
