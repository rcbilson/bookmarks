import React from "react";

import NavWidget from "./NavWidget.tsx";
import BookmarkQuery from "./BookmarkQuery.tsx";

const FavoritePage: React.FC = () => {
  return (
    <div id="recentContainer">
      <NavWidget/>
      <BookmarkQuery queryPath='/api/favorites?count=10' />
    </div>
  );
};

export default FavoritePage;
