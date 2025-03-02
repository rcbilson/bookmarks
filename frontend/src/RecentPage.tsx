// A react component that has an editable text area for a recipe url
// next to a button with a refresh icon. When the button is clicked,
// the recipe url is fetched and the text area below the url is updated
// with the recipe contents.
import React from "react";

import NavWidget from "./NavWidget.tsx";
import BookmarkQuery from "./BookmarkQuery.tsx";

const RecentPage: React.FC = () => {
  return (
    <div id="recentContainer">
      <NavWidget/>
      <BookmarkQuery queryPath='/api/recents?count=10' />
    </div>
  );
};

export default RecentPage;
