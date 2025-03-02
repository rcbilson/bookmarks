import React from "react";
import axios from "axios";
import { useNavigate } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query'

interface Props {
  contents?: string;
}

const NavWidget: React.FC<Props> = ({contents}: Props) => {
  const navigate = useNavigate();
  const queryClient = useQueryClient()

  const handleSearchTextChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const searchText = event.target.value;
    if (!searchText) {
      navigate("/");
      return;
    }
    try {
      new URL(searchText);
      await axios.post("/api/add?url=" + encodeURIComponent(searchText));
      queryClient.invalidateQueries({ queryKey: ['bookmarkList'] })
      navigate("/recent");
    } catch (_) {
      navigate("/search?q=" + encodeURIComponent(searchText));
    }
  };

  return (
    <div id="searchbar">
      <input id="url" type="text" value={contents} onChange={handleSearchTextChange} />
      <div id="navlinks">
        <a id="recentlink" href="/recent">Recent</a>
        <a id="favoritelink" href="/favorite">Favorites</a>
      </div>
    </div>
  )
};

export default NavWidget;
