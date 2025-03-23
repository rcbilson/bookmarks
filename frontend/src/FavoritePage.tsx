import React from "react";
import BookmarkQuery from "./BookmarkQuery";

const FavoritePage: React.FC = () => {
  return <BookmarkQuery queryPath='/api/favorites?count=10' />;
};

export default FavoritePage;
