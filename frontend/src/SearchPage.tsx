import React, { useState, useEffect } from "react";
import { Input } from '@chakra-ui/react';

import BookmarkQuery from "./BookmarkQuery.tsx";

const SearchPage: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery);
    }, 500);

    return () => clearTimeout(timer);
  }, [searchQuery]);

  return (
    <div>
      <Input
        placeholder="Search recipes..."
        value={searchQuery}
        onChange={(e) => setSearchQuery(e.target.value)}
        mb={4}
      />
      {debouncedQuery && <BookmarkQuery queryPath={"/api/search?q=" + encodeURIComponent(debouncedQuery)} />}
    </div>
  );
};

export default SearchPage;